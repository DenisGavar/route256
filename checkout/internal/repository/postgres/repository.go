package repository

import (
	"context"
	"route256/checkout/internal/domain/model"
	"route256/libs/transactor"
)

const (
	basketsTable = "baskets"
)

type CheckoutRepository interface {
	AddToCart(ctx context.Context, addToCartRequest *model.AddToCartRequest) error
	ListCart(ctx context.Context, listCartRequest *model.ListCartRequest) (*model.ListCartResponse, error)
	DeleteFromCart(ctx context.Context, deleteFromCartRequest *model.DeleteFromCartRequest) error
}

var _ CheckoutRepository = (*repository)(nil)

type repository struct {
	queryEngineProvider transactor.QueryEngineProvider
}

func NewRepository(queryEngineProvider transactor.QueryEngineProvider) *repository {
	return &repository{
		queryEngineProvider: queryEngineProvider,
	}
}
