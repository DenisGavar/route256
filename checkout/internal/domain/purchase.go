package domain

import (
	"context"
)

type PurchaseRequest struct {
	// user ID
	User int64
}

type PurchaseResponse struct {
	OrderId int64
}

type Order struct {
	OrderID int64
}

func (m *model) Purchase(ctx context.Context, req *PurchaseRequest) (*PurchaseResponse, error) {

	return &PurchaseResponse{
		OrderId: 42,
	}, nil

	// order, err := m.orderCreator.CreateOrder(ctx, user)
	// if err != nil {
	// 	return order, errors.WithMessage(err, "create order")
	// }

	// return order, nil
}
