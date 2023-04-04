package loms_v1

import (
	"context"
	"fmt"
	"route256/libs/logger"
	"route256/loms/internal/converter"
	desc "route256/loms/pkg/loms_v1"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) CancelOrder(ctx context.Context, req *desc.CancelOrderRequest) (*emptypb.Empty, error) {
	logger.Debug("loms server", zap.String("handler", "CancelOrder"), zap.String("request", fmt.Sprintf("%+v", req)))

	span, ctx := opentracing.StartSpanFromContext(ctx, "CancelOrder processing")
	defer span.Finish()

	span.SetTag("order_id", req.GetOrderId())

	err := i.lomsService.CancelOrder(ctx, converter.FromDescToMolelCancelOrderRequest(req))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
