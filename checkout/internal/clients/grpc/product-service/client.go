package product_service

import (
	"context"
	product "route256/checkout/pkg/product-service_v1"

	"google.golang.org/grpc"
)

//go:generate mockgen -source="client.go" -destination="mocks/client_mock.go" -package=mocks . ProductServiceClient
type ProductServiceClient interface {
	GetProduct(context.Context, *product.GetProductRequest) (*product.GetProductResponse, error)
}

type client struct {
	productServiceClient product.ProductServiceClient

	token string
}

func New(cc *grpc.ClientConn, token string) *client {
	return &client{
		productServiceClient: product.NewProductServiceClient(cc),
		token:                token,
	}
}
