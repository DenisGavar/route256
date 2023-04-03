package loms_v1

import (
	"context"
	"fmt"
	"route256/libs/logger"
	"route256/loms/internal/converter"
	desc "route256/loms/pkg/loms_v1"

	"go.uber.org/zap"
)

func (i *Implementation) ListOrder(ctx context.Context, req *desc.ListOrderRequest) (*desc.ListOrderResponse, error) {
	logger.Debug("loms server", zap.String("handler", "ListOrder"), zap.String("request", fmt.Sprintf("%+v", req)))

	response, err := i.lomsService.ListOrder(ctx, converter.FromDescToMolelListOrderRequest(req))
	if err != nil {
		return nil, err
	}

	return converter.FromModelToDescListOrderResponse(response), nil
}
