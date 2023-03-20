package checkout_v1

import (
	"context"
	"log"
	"route256/checkout/internal/converter"
	desc "route256/checkout/pkg/checkout_v1"
)

func (i *Implementation) Purchase(ctx context.Context, req *desc.PurchaseRequest) (*desc.PurchaseResponse, error) {
	log.Printf("purchase: %+v", req)

	response, err := i.checkoutModel.Purchase(ctx, converter.FromDescToMolelPurchaseRequest(req))
	if err != nil {
		return nil, err
	}

	return converter.FromMolelToDescPurchaseResponse(response), nil
}
