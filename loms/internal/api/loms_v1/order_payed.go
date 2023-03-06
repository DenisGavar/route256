package loms_v1

import (
	"context"
	"log"
	"route256/loms/internal/converter"
	desc "route256/loms/pkg/loms_v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) OrderPayed(ctx context.Context, req *desc.OrderPayedRequest) (*emptypb.Empty, error) {
	log.Printf("order payed: %+v", req)

	err := i.lomsService.OrderPayed(ctx, converter.FromDescToMolelOrderPayedRequest(req))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
