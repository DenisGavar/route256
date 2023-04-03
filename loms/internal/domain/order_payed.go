package domain

import (
	"context"
	"fmt"
	"route256/libs/logger"
	"route256/loms/internal/domain/model"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (s *service) OrderPayed(ctx context.Context, req *model.OrderPayedRequest) error {
	// оплачиваем заказ
	logger.Debug("loms domain", zap.String("handler", "OrderPayed"), zap.String("request", fmt.Sprintf("%+v", req)))

	// очищаем резервы
	err := s.repository.lomsRepository.ClearReserves(ctx, req.OrderId)
	if err != nil {
		return errors.WithMessage(err, ErrClearingReserves.Error())
	}

	// вызываем метод смены статуса
	// payed
	err = s.repository.lomsRepository.ChangeStatus(ctx, req.OrderId, model.OrderStatusPayed)
	if err != nil {
		return errors.WithMessage(err, ErrChangingStatus.Error())
	}

	return nil
}
