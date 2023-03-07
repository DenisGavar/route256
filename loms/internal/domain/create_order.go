package domain

import (
	"context"
	"route256/loms/internal/domain/model"
)

func (m *service) CreateOrder(ctx context.Context, req *model.CreateOrderRequest) (*model.CreateOrderResponse, error) {
	var response *model.CreateOrderResponse
	err := m.repository.transactionManager.RunRepeatableRead(ctx, func(ctxTX context.Context) error {
		// создаём заказ, получаем его id
		repoResponse, err := m.repository.lomsRepository.CreateOrder(ctx, req)
		if err != nil {
			return err
		}
		response = repoResponse

		// дополняем структуру orderID
		req.OrderId = response.OrderId

		// резервируем товары из заказа с его id
		err = m.repository.lomsRepository.ReserveItems(ctx, req)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}
