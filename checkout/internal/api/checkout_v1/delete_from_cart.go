package checkout_v1

import (
	"context"
	"log"
	desc "route256/checkout/pkg/checkout_v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) DeleteFromCart(ctx context.Context, req *desc.DeleteFromCartRequest) (*emptypb.Empty, error) {
	log.Printf("deleteFromCart: %+v", req)

	return &emptypb.Empty{}, nil
}
