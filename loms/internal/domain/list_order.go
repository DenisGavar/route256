package domain

import (
	"context"
	"route256/loms/internal/domain/model"
)

func (m *service) ListOrder(ctx context.Context, req *model.ListOrderRequest) (*model.ListOrderResponse, error) {
	return &model.ListOrderResponse{
		Status: model.OderStatus_awaiting_payment,
		User:   15,
		Items: []*model.OrderItem{
			{
				Sku:   33,
				Count: 2,
			},
			{
				Sku:   44,
				Count: 3,
			},
		},
	}, nil
}
