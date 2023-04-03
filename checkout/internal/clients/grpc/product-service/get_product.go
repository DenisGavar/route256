package product_service

import (
	"context"
	"fmt"
	product "route256/checkout/pkg/product-service_v1"
	"route256/libs/logger"

	"go.uber.org/zap"
)

func (c *client) GetProduct(ctx context.Context, req *product.GetProductRequest) (*product.GetProductResponse, error) {
	logger.Debug("product-service client", zap.String("handler", "GetProduct"), zap.String("request", fmt.Sprintf("%+v", req)))

	req.Token = c.token

	response, err := c.productServiceClient.GetProduct(ctx, req)
	if err != nil {
		return nil, err
	}

	return response, nil
}
