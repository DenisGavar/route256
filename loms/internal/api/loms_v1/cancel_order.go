package loms_v1

import (
	"context"
	"log"
	"route256/loms/internal/converter"
	desc "route256/loms/pkg/loms_v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) CancelOrder(ctx context.Context, req *desc.CancelOrderRequest) (*emptypb.Empty, error) {
	log.Printf("cancel order: %+v", req)

	err := i.lomsService.CancelOrder(ctx, converter.FromDescToMolelCancelOrderRequest(req))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
