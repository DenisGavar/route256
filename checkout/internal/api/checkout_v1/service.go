package checkout_v1

import (
	desc "route256/checkout/pkg/checkout_v1"
)

type Implementation struct {
	desc.UnimplementedCheckoutV1Server
}

func NewCheckoutV1() *Implementation {
	return &Implementation{
		desc.UnimplementedCheckoutV1Server{},
	}
}