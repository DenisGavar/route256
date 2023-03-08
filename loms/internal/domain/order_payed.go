package domain

import (
	"context"
	"route256/loms/internal/domain/model"
)

func (m *service) OrderPayed(ctx context.Context, req *model.OrderPayedRequest) error {
	// очищаем резервы
	err := m.repository.lomsRepository.ClearReserves(ctx, req.OrderId)
	if err != nil {
		return err
	}

	// вызываем метод смены статуса
	// payed
	err = m.repository.lomsRepository.ChangeStatus(ctx, req.OrderId, model.OrderStatusPayed)
	if err != nil {
		return err
	}

	return nil
}
