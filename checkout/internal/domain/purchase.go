package domain

import (
	"context"
	loms "route256/loms/pkg/loms_v1"

	"github.com/pkg/errors"
)

type PurchaseRequest struct {
	// user ID
	User int64
}

type PurchaseResponse struct {
	OrderId int64
}

func (m *model) Purchase(ctx context.Context, req *PurchaseRequest) (*PurchaseResponse, error) {

	order, err := m.lomsClient.CreateOrder(ctx, &loms.CreateOrderRequest{User: req.User})
	if err != nil {
		return nil, errors.WithMessage(err, "create order")
	}

	return &PurchaseResponse{OrderId: order.GetOrderId()}, nil
}
