package checkout_v1

import (
	"context"
	desc "route256/checkout/pkg/checkout_v1"
)

func (i *Implementation) ListCart(ctx context.Context, req *desc.ListCartRequest) (*desc.ListCartResponse, error) {
	return &desc.ListCartResponse{
		Items: []*desc.CartItem{
			{
				Sku:   111,
				Count: 1,
				Name:  "qwe",
				Price: 1,
			},
			{
				Sku:   222,
				Count: 2,
				Name:  "qwe",
				Price: 2,
			},
		},
		TotalPrice: 5,
	}, nil
}
