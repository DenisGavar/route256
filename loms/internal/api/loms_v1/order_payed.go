package loms_v1

import (
	"context"
	"log"
	desc "route256/loms/pkg/loms_v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) OrderPayed(ctx context.Context, req *desc.OrderPayedRequest) (*emptypb.Empty, error) {
	log.Printf("order payed: %+v", req)

	return &emptypb.Empty{}, nil
}