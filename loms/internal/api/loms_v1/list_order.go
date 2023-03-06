package loms_v1

import (
	"context"
	"log"
	"route256/loms/internal/converter"
	desc "route256/loms/pkg/loms_v1"
)

func (i *Implementation) ListOrder(ctx context.Context, req *desc.ListOrderRequest) (*desc.ListOrderResponse, error) {
	log.Printf("list order: %+v", req)

	response, err := i.lomsService.ListOrder(ctx, converter.FromDescToMolelListOrderRequest(req))
	if err != nil {
		return nil, err
	}

	return converter.FromModelToDescListOrderResponse(response), nil
}
