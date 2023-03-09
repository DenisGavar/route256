package domain

import (
	"context"
	"route256/checkout/internal/domain/model"
	loms "route256/loms/pkg/loms_v1"

	"github.com/pkg/errors"
)

func (s *service) Purchase(ctx context.Context, req *model.PurchaseRequest) (*model.PurchaseResponse, error) {

	order, err := s.lomsClient.CreateOrder(ctx, &loms.CreateOrderRequest{User: req.User})
	if err != nil {
		return nil, errors.WithMessage(err, "create order")
	}

	return &model.PurchaseResponse{OrderId: order.GetOrderId()}, nil
}
