package domain

import (
	"context"
	"fmt"
	"route256/libs/logger"
	"route256/loms/internal/domain/model"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (s *service) CreateOrder(ctx context.Context, req *model.CreateOrderRequest) (*model.CreateOrderResponse, error) {
	// создаём заказ, получаем его id
	logger.Debug("loms domain", zap.String("handler", "CreateOrder"), zap.String("request", fmt.Sprintf("%+v", req)))

	response, err := s.repository.lomsRepository.CreateOrder(ctx, req)
	if err != nil {
		return nil, errors.WithMessage(err, ErrCreatingOrder.Error())
	}

	// дополняем структуру orderID
	req.OrderId = response.OrderId

	err = s.repository.transactionManager.RunRepeatableRead(ctx, func(ctxTX context.Context) error {
		// резервируем товары из заказа с его id
		for _, orderItem := range req.Items {
			// проверяем наличие каждого товара на складах
			stocks, err := s.repository.lomsRepository.Stocks(ctxTX, &model.StocksRequest{Sku: orderItem.Sku})
			if err != nil {
				return errors.WithMessage(err, ErrGettingStocks.Error())
			}

			var reservedCount uint64
			// собираем на каких складах сколько нужно взять товара
			needToReserve := make([]*model.ReserveStocksItem, 0, len(stocks.Stocks))

			for _, stocksItem := range stocks.Stocks {

				left := uint64(orderItem.Count) - reservedCount
				if left == 0 {
					break
				}

				if stocksItem.Count > left {
					needToReserve = append(needToReserve, &model.ReserveStocksItem{
						WarehouseId: stocksItem.WarehouseId,
						Sku:         orderItem.Sku,
						Count:       left,
					})
					reservedCount += left
				} else {
					needToReserve = append(needToReserve, &model.ReserveStocksItem{
						WarehouseId: stocksItem.WarehouseId,
						Sku:         orderItem.Sku,
						Count:       stocksItem.Count,
					})
					reservedCount += stocksItem.Count
				}
			}

			if reservedCount != uint64(orderItem.Count) {
				return ErrNotEnoughItems
			}

			// резервируем товары
			for _, reserveStockItem := range needToReserve {
				if err := s.repository.lomsRepository.ReserveItems(ctxTX, response.OrderId, reserveStockItem); err != nil {
					return errors.WithMessage(err, ErrReservingItems.Error())
				}
			}
		}

		return nil
	})

	// проверяем успешность резерва
	if err != nil {
		// вызываем метод смены статуса
		// failed
		err = s.repository.lomsRepository.ChangeStatus(ctx, response.OrderId, model.OrderStatusFailed)
		if err != nil {
			return nil, errors.WithMessage(err, ErrChangingStatus.Error())
		}
	} else {
		// вызываем метод смены статуса
		// awaiting payment
		err = s.repository.lomsRepository.ChangeStatus(ctx, response.OrderId, model.OrderStatusAwaitingPayment)
		if err != nil {
			return nil, errors.WithMessage(err, ErrChangingStatus.Error())
		}
	}

	return response, nil
}
