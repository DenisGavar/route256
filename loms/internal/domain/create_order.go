package domain

import (
	"context"
	"route256/loms/internal/domain/model"
)

func (m *service) CreateOrder(ctx context.Context, req *model.CreateOrderRequest) (*model.CreateOrderResponse, error) {
	return &model.CreateOrderResponse{OrderId: 42}, nil
}
