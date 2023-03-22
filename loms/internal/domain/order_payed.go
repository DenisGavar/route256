package domain

import (
	"context"
	"route256/loms/internal/domain/model"

	"github.com/pkg/errors"
)

func (s *service) OrderPayed(ctx context.Context, req *model.OrderPayedRequest) error {
	// оплачиваем заказ

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
