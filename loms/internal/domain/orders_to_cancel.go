package domain

import (
	"context"
	"route256/libs/logger"
	"route256/loms/internal/domain/model"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (s *service) OrdersToCancel(ctx context.Context, time time.Time) ([]*model.CancelOrderRequest, error) {
	// получаем заказы на отмену
	logger.Debug("loms domain", zap.String("handler", "OrdersToCancel"), zap.String("time", time.String()))

	response, err := s.repository.lomsRepository.OrdersToCancel(ctx, time)
	if err != nil {
		return nil, errors.WithMessage(err, ErrGettingOrdersToCancel.Error())
	}

	return response, nil
}
