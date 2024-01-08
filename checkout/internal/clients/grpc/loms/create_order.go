package loms

import (
	"context"
	"fmt"
	"route256/libs/logger"
	loms "route256/loms/pkg/loms_v1"

	"go.uber.org/zap"
)

func (c *client) CreateOrder(ctx context.Context, req *loms.CreateOrderRequest) (*loms.CreateOrderResponse, error) {
	logger.Debug("loms client", zap.String("handler", "CreateOrder"), zap.String("request", fmt.Sprintf("%+v", req)))

	response, err := c.lomsClient.CreateOrder(ctx, req)
	if err != nil {
		return nil, err
	}

	return response, nil
}
