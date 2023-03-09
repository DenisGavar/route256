package checkout_v1

import (
	"context"
	"log"
	"route256/checkout/internal/converter"
	desc "route256/checkout/pkg/checkout_v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) DeleteFromCart(ctx context.Context, req *desc.DeleteFromCartRequest) (*emptypb.Empty, error) {
	log.Printf("deleteFromCart: %+v", req)

	err := i.checkoutModel.DeleteFromCart(ctx, converter.FromDescToMolelDeleteFromCartRequest(req))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
