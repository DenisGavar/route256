package checkout_v1

import (
	"context"
	"fmt"
	"route256/checkout/internal/converter"
	desc "route256/checkout/pkg/checkout_v1"
	"route256/libs/logger"

	"go.uber.org/zap"
)

func (i *Implementation) ListCart(ctx context.Context, req *desc.ListCartRequest) (*desc.ListCartResponse, error) {
	logger.Debug("checkout server", zap.String("handler", "ListCart"), zap.String("request", fmt.Sprintf("%+v", req)))

	response, err := i.checkoutModel.ListCart(ctx, converter.FromDescToMolelListCartRequest(req))
	if err != nil {
		return nil, err
	}

	return converter.FromMolelToDescListCartResponse(response), nil
}
