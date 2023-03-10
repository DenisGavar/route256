package checkout_v1

import (
	"route256/checkout/internal/domain"
	desc "route256/checkout/pkg/checkout_v1"
)

type Implementation struct {
	desc.UnimplementedCheckoutV1Server

	checkoutModel domain.Model
}

func NewCheckoutV1(checkoutModel domain.Model) *Implementation {
	return &Implementation{
		desc.UnimplementedCheckoutV1Server{},
		checkoutModel,
	}
}
