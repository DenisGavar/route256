package product_service

import (
	"context"
	product "route256/checkout/pkg/product-service_v1"
)

func (c *client) GetProduct(ctx context.Context, req *product.GetProductRequest) (*product.GetProductResponse, error) {
	req.Token = c.token

	response, err := c.productServiceClient.GetProduct(ctx, req)
	if err != nil {
		return nil, err
	}

	return response, nil
}
