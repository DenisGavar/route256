package domain

import (
	"context"
	product "route256/checkout/pkg/product-service_v1"
	loms "route256/loms/pkg/loms_v1"
)

var _ Model = (*model)(nil)

type Model interface {
	AddToCart(context.Context, *AddToCartRequest) error
	ListCart(context.Context, *ListCartRequest) (*ListCartResponse, error)
	Purchase(context.Context, *PurchaseRequest) (*PurchaseResponse, error)
}

type ProductServiceClient interface {
	GetProduct(context.Context, *product.GetProductRequest) (*product.GetProductResponse, error)
}

type LomsClient interface {
	Stocks(context.Context, *loms.StocksRequest) (*loms.StocksResponse, error)
	CreateOrder(context.Context, *loms.CreateOrderRequest) (*loms.CreateOrderResponse, error)
}

type model struct {
	lomsClient           LomsClient
	productServiceClient ProductServiceClient
}

func NewModel(lomsClient LomsClient, productServiceClient ProductServiceClient) *model {
	return &model{
		lomsClient:           lomsClient,
		productServiceClient: productServiceClient,
	}
}
