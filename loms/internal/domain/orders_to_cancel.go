package domain

import (
	"context"
	"route256/loms/internal/domain/model"
	"time"

	"github.com/pkg/errors"
)

func (s *service) OrdersToCancel(ctx context.Context, time time.Time) ([]*model.CancelOrderRequest, error) {
	// получаем заказы на отмену

	response, err := s.repository.lomsRepository.OrdersToCancel(ctx, time)
	if err != nil {
		return nil, errors.WithMessage(err, "getting orders to cancel")
	}

	return response, nil
}
