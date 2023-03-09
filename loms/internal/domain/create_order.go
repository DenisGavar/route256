package domain

import (
	"context"
	"route256/loms/internal/domain/model"
)

func (s *service) CreateOrder(ctx context.Context, req *model.CreateOrderRequest) (*model.CreateOrderResponse, error) {
	// создаём заказ, получаем его id
	response, err := s.repository.lomsRepository.CreateOrder(ctx, req)
	if err != nil {
		return nil, err
	}

	// дополняем структуру orderID
	req.OrderId = response.OrderId

	err = s.repository.transactionManager.RunRepeatableRead(ctx, func(ctxTX context.Context) error {
		// резервируем товары из заказа с его id
		for _, orderItem := range req.Items {
			// проверяем наличие каждого товара на складах
			stocks, err := s.repository.lomsRepository.Stocks(ctxTX, &model.StocksRequest{Sku: orderItem.Sku})
			if err != nil {
				return err
			}

			var reservedCount uint64
			// собираем на каких складах нсколько нужно взять товара
			needToReserve := make(map[int64]model.ReserveStocksItem, 1)

			for _, stocksItem := range stocks.Stocks {

				left := uint64(orderItem.Count) - reservedCount
				if left == 0 {
					break
				}

				if stocksItem.Count > left {
					needToReserve[stocksItem.WarehouseId] = model.ReserveStocksItem{
						Sku:   orderItem.Sku,
						Count: left,
						Part:  true,
					}
					reservedCount += left
				} else {
					needToReserve[stocksItem.WarehouseId] = model.ReserveStocksItem{
						Sku:   orderItem.Sku,
						Count: stocksItem.Count,
						Part:  false,
					}
					reservedCount += stocksItem.Count
				}
			}

			if reservedCount != uint64(orderItem.Count) {
				return ErrNotEnoughItems
			}

			// резервируем товары
			for warehouseId, reserveStockItem := range needToReserve {
				if err := s.repository.lomsRepository.ReserveItems(ctxTX, response.OrderId, warehouseId, &reserveStockItem); err != nil {
					return err
				}
			}
		}

		return nil
	})

	// проверяем успешность резерва
	if err != nil {
		// тут немного странно, т.к. получили ошибку, но потом её можем проигнорировать
		// пока так, т.к. потом мы это поменяем, насколько я понял
		// будем создавать заказ и сразу возвращать его ID, а резервировать будем в отдельной горутине

		// вызываем метод смены статуса
		// failed
		err = s.repository.lomsRepository.ChangeStatus(ctx, response.OrderId, model.OrderStatusFailed)
		if err != nil {
			return nil, err
		}
	} else {
		// вызываем метод смены статуса
		// awaiting payment
		err = s.repository.lomsRepository.ChangeStatus(ctx, response.OrderId, model.OrderStatusAwaitingPayment)
		if err != nil {
			return nil, err
		}
	}

	return response, nil
}
