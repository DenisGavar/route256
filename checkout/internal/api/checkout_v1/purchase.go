package checkout_v1

import (
	"context"
	"fmt"
	"route256/checkout/internal/converter"
	desc "route256/checkout/pkg/checkout_v1"
	"route256/libs/logger"

	"go.uber.org/zap"
)

func (i *Implementation) Purchase(ctx context.Context, req *desc.PurchaseRequest) (*desc.PurchaseResponse, error) {
	logger.Debug("checkout server", zap.String("handler", "Purchase"), zap.String("request", fmt.Sprintf("%+v", req)))

	response, err := i.checkoutModel.Purchase(ctx, converter.FromDescToMolelPurchaseRequest(req))
	if err != nil {
		return nil, err
	}

	return converter.FromMolelToDescPurchaseResponse(response), nil
}
