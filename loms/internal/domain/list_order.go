package domain

import (
	"context"
	"route256/loms/internal/domain/model"
)

func (m *service) ListOrder(ctx context.Context, req *model.ListOrderRequest) (*model.ListOrderResponse, error) {
	response, err := m.repository.lomsRepository.ListOrder(ctx, req)
	if err != nil {
		return nil, err
	}

	// нужна проверка, если ничего не вернулось
	// либо если вернулся заказ без строк

	return response, nil
}
