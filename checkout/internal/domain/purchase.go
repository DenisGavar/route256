package domain

import (
	"context"

	"github.com/pkg/errors"
)

type Order struct {
	OrderID int64
}

func (m *Model) Purchase(ctx context.Context, user int64) (*Order, error) {

	order, err := m.orderCreator.CreateOrder(ctx, user)
	if err != nil {
		return order, errors.WithMessage(err, "create order")
	}

	return order, nil
}
