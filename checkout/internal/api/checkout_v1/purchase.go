package checkout_v1

import (
	"context"
	desc "route256/checkout/pkg/checkout_v1"
)

func (i *Implementation) Purchase(ctx context.Context, req *desc.PurchaseRequest) (*desc.PurchaseResponse, error) {
	return &desc.PurchaseResponse{
		OrderId: 42,
	}, nil
}
