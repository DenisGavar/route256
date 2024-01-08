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

func (i *Implementation) Stocks(ctx context.Context, req *desc.StocksRequest) (*desc.StocksResponse, error) {
	logger.Debug("loms server", zap.String("handler", "Stocks"), zap.String("request", fmt.Sprintf("%+v", req)))

	span, ctx := opentracing.StartSpanFromContext(ctx, "Stocks processing")
	defer span.Finish()

	span.SetTag("sku", req.GetSku())

	response, err := i.lomsService.Stocks(ctx, converter.FromDescToMolelStocksRequest(req))
	if err != nil {
		return nil, err
	}

	return converter.FromModelToDescStocksResponse(response), nil
}
