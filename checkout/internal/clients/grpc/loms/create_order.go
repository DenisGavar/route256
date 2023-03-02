package loms

import (
	"context"
	loms "route256/loms/pkg/loms_v1"
)

func (c *client) CreateOrder(ctx context.Context, request *loms.CreateOrderRequest) (*loms.CreateOrderResponse, error) {
	response, err := c.lomsClient.CreateOrder(ctx, request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
