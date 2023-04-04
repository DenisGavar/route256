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

func (i *Implementation) AddToCart(ctx context.Context, req *desc.AddToCartRequest) (*emptypb.Empty, error) {
	logger.Debug("checkout server", zap.String("handler", "AddToCart"), zap.String("request", fmt.Sprintf("%+v", req)))

	// нулевое количество добавлять нет смысла
	if req.Count == 0 {
		return nil, ErrNullCount
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "AddToCart processing")
	defer span.Finish()

	span.SetTag("user", req.GetUser())
	span.SetTag("sku", req.GetSku())

	err := i.checkoutModel.AddToCart(ctx, converter.FromDescToMolelAddToCartRequest(req))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
