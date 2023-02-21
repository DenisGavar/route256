package domain

import (
	"context"

	"github.com/pkg/errors"
)

func (m *Model) Purchase(ctx context.Context, user int64) (int64, error) {

	orderID, err := m.orderCreator.CreateOrder(ctx, user)
	if err != nil {
		return 0, errors.WithMessage(err, "create order")
	}

	return orderID, nil
}
