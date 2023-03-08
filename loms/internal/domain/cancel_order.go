package domain

import (
	"context"
	"route256/loms/internal/domain/model"
)

func (m *service) CancelOrder(ctx context.Context, req *model.CancelOrderRequest) error {

	err := m.repository.transactionManager.RunRepeatableRead(ctx, func(ctxTX context.Context) error {
		// получаем резервы, которые надо вернуть на склад

		// возвращаем резервы на склад

		// очищаем резервы
		err := m.repository.lomsRepository.ClearReserves(ctx, req.OrderId)
		if err != nil {
			return err
		}

		// вызываем метод смены статуса
		// failed
		err = m.repository.lomsRepository.ChangeStatus(ctx, req.OrderId, model.OrderStatusCancelled)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
