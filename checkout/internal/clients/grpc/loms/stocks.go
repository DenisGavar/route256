package loms

import (
	"context"
	loms "route256/loms/pkg/loms_v1"
)

func (c *client) Stocks(ctx context.Context, request *loms.StocksRequest) (*loms.StocksResponse, error) {
	response, err := c.lomsClient.Stocks(ctx, request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
