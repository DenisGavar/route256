package loms_v1

import (
	"context"
	"log"
	desc "route256/loms/pkg/loms_v1"
)

func (i *Implementation) CreateOrder(ctx context.Context, req *desc.CreateOrderRequest) (*desc.CreateOrderResponse, error) {
	log.Printf("create order: %+v", req)

	return &desc.CreateOrderResponse{
		OrderId: 42,
	}, nil
}
