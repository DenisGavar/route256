package product_service

import (
	product "route256/checkout/pkg/product-service_v1"

	"google.golang.org/grpc"
)

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
