package domain

import (
	"context"
	"route256/loms/internal/domain/model"
)

func (m *service) CreateOrder(ctx context.Context, req *model.CreateOrderRequest) (*model.CreateOrderResponse, error) {
	// Создает новый заказ для пользователя из списка переданных товаров. Товары при этом нужно зарезервировать на складе.

	response, err := m.repository.lomsRepository.CreateOrder(ctx, req)
	if err != nil {
		return nil, err
	}

	return response, nil
}
