package loms_v1

import (
	"context"
	desc "route256/loms/pkg/loms_v1"
)

func (i *Implementation) CreateOrder(ctx context.Context, req *desc.CreateOrderRequest) (*desc.CreateOrderResponse, error) {
	return &desc.CreateOrderResponse{
		OrderId: 42,
	}, nil
}
