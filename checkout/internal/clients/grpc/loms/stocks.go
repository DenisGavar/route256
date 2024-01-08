package loms

import (
	"context"
	"fmt"
	"route256/libs/logger"
	loms "route256/loms/pkg/loms_v1"

	"go.uber.org/zap"
)

func (c *client) Stocks(ctx context.Context, req *loms.StocksRequest) (*loms.StocksResponse, error) {
	logger.Debug("loms client", zap.String("handler", "Stocks"), zap.String("request", fmt.Sprintf("%+v", req)))

	response, err := c.lomsClient.Stocks(ctx, req)
	if err != nil {
		return nil, err
	}

	return response, nil
}
