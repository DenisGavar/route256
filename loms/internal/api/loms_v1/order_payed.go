package loms_v1

import (
	"context"
	"fmt"
	"route256/libs/logger"
	"route256/loms/internal/converter"
	desc "route256/loms/pkg/loms_v1"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) OrderPayed(ctx context.Context, req *desc.OrderPayedRequest) (*emptypb.Empty, error) {
	logger.Debug("loms server", zap.String("handler", "OrderPayed"), zap.String("request", fmt.Sprintf("%+v", req)))

	err := i.lomsService.OrderPayed(ctx, converter.FromDescToMolelOrderPayedRequest(req))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
