package domain

import (
	"context"
	"route256/loms/internal/domain/model"
)

func (s *service) OrderPayed(ctx context.Context, req *model.OrderPayedRequest) error {
	// очищаем резервы
	err := s.repository.lomsRepository.ClearReserves(ctx, req.OrderId)
	if err != nil {
		return err
	}

	// вызываем метод смены статуса
	// payed
	err = s.repository.lomsRepository.ChangeStatus(ctx, req.OrderId, model.OrderStatusPayed)
	if err != nil {
		return err
	}

	return nil
}
