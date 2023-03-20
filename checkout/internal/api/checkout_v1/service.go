package checkout_v1

import (
	"route256/checkout/internal/domain"
	desc "route256/checkout/pkg/checkout_v1"

	"github.com/pkg/errors"
)

var (
	ErrNullCount = errors.New("null count is not allowed")
)

type Implementation struct {
	desc.UnimplementedCheckoutV1Server

	checkoutModel domain.Service
}

func NewCheckoutV1(checkoutModel domain.Service) *Implementation {
	return &Implementation{
		desc.UnimplementedCheckoutV1Server{},
		checkoutModel,
	}
}
