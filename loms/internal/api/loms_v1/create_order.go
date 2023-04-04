package loms_v1

import (
	"context"
	"fmt"
	"route256/libs/logger"
	"route256/loms/internal/converter"
	desc "route256/loms/pkg/loms_v1"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

func (i *Implementation) CreateOrder(ctx context.Context, req *desc.CreateOrderRequest) (*desc.CreateOrderResponse, error) {
	logger.Debug("loms server", zap.String("handler", "CreateOrder"), zap.String("request", fmt.Sprintf("%+v", req)))

	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateOrder processing")
	defer span.Finish()

	span.SetTag("user", req.GetUser())

	response, err := i.lomsService.CreateOrder(ctx, converter.FromDescToMolelCreateOrderRequest(req))
	if err != nil {
		return nil, err
	}

	return converter.FromModelToDescCreateOrderResponse(response), nil
}
