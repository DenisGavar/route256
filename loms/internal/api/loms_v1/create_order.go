package loms_v1

import (
	"context"
	"log"
	"route256/loms/internal/converter"
	desc "route256/loms/pkg/loms_v1"
)

func (i *Implementation) CreateOrder(ctx context.Context, req *desc.CreateOrderRequest) (*desc.CreateOrderResponse, error) {
	log.Printf("create order: %+v", req)

	response, err := i.lomsService.CreateOrder(ctx, converter.FromDescToMolelCreateOrderRequest(req))
	if err != nil {
		return nil, err
	}

	return converter.FromModelToDescCreateOrderResponse(response), nil
}
