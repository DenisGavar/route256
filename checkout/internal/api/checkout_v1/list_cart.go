package checkout_v1

import (
	"context"
	"log"
	"route256/checkout/internal/converter"
	desc "route256/checkout/pkg/checkout_v1"
)

func (i *Implementation) ListCart(ctx context.Context, req *desc.ListCartRequest) (*desc.ListCartResponse, error) {
	log.Printf("listCart: %+v", req)

	response, err := i.checkoutModel.ListCart(ctx, converter.FromDescToMolelListCartRequest(req))
	if err != nil {
		return nil, err
	}

	return converter.FromMolelToDescListCartResponse(response), nil
}
