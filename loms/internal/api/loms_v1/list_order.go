package loms_v1

import (
	"context"
	"log"
	desc "route256/loms/pkg/loms_v1"
)

func (i *Implementation) ListOrder(ctx context.Context, req *desc.ListOrderRequest) (*desc.ListOrderResponse, error) {
	log.Printf("list order: %+v", req)

	return &desc.ListOrderResponse{
		Status: desc.OderStatus_awaiting_payment,
		User:   15,
		Items: []*desc.OrderItem{
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
