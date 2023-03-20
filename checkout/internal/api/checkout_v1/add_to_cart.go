package checkout_v1

import (
	"context"
	"log"
	"route256/checkout/internal/converter"
	desc "route256/checkout/pkg/checkout_v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) AddToCart(ctx context.Context, req *desc.AddToCartRequest) (*emptypb.Empty, error) {
	log.Printf("addToCart: %+v", req)

	// нулевое количество добавлять нет смысла
	if req.Count == 0 {
		return nil, ErrNullCount
	}

	err := i.checkoutModel.AddToCart(ctx, converter.FromDescToMolelAddToCartRequest(req))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
