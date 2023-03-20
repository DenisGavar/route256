package domain

import (
	"context"
	"route256/loms/internal/domain/model"

	"github.com/pkg/errors"
)

func (s *service) CreateOrder(ctx context.Context, req *model.CreateOrderRequest) (*model.CreateOrderResponse, error) {
	// создаём заказ, получаем его id
	response, err := s.repository.lomsRepository.CreateOrder(ctx, req)
	if err != nil {
		return nil, errors.WithMessage(err, "creating order")
	}

	// дополняем структуру orderID
	req.OrderId = response.OrderId

	err = s.repository.transactionManager.RunRepeatableRead(ctx, func(ctxTX context.Context) error {
		// резервируем товары из заказа с его id
		for _, orderItem := range req.Items {
			// проверяем наличие каждого товара на складах
			stocks, err := s.repository.lomsRepository.Stocks(ctxTX, &model.StocksRequest{Sku: orderItem.Sku})
			if err != nil {
				return errors.WithMessage(err, "checking stocks")
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
					return errors.WithMessage(err, "reserving items")
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
			return nil, errors.WithMessage(err, "changing ctatus")
		}
	} else {
		// вызываем метод смены статуса
		// awaiting payment
		err = s.repository.lomsRepository.ChangeStatus(ctx, response.OrderId, model.OrderStatusAwaitingPayment)
		if err != nil {
			return nil, errors.WithMessage(err, "changing ctatus")
		}
	}

	return response, nil
}
