package checkout_v1

import (
	"context"
	"fmt"
	"route256/checkout/internal/converter"
	desc "route256/checkout/pkg/checkout_v1"
	"route256/libs/logger"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) DeleteFromCart(ctx context.Context, req *desc.DeleteFromCartRequest) (*emptypb.Empty, error) {
	logger.Debug("checkout server", zap.String("handler", "DeleteFromCart"), zap.String("request", fmt.Sprintf("%+v", req)))

	// нулевое количество удалять нет смысла
	if req.Count == 0 {
		return nil, ErrNullCount
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "DeleteFromCart processing")
	defer span.Finish()

	span.SetTag("user", req.GetUser())
	span.SetTag("sku", req.GetSku())

	err := i.checkoutModel.DeleteFromCart(ctx, converter.FromDescToMolelDeleteFromCartRequest(req))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
