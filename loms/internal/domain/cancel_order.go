package domain

import (
	"context"
	"route256/loms/internal/domain/model"
)

func (m *service) CancelOrder(ctx context.Context, req *model.CancelOrderRequest) error {

	err := m.repository.transactionManager.RunRepeatableRead(ctx, func(ctxTX context.Context) error {
		// получаем резервы, которые надо вернуть на склад
		needToDropReserve, err := m.repository.lomsRepository.Reserves(ctx, req.OrderId)
		if err != nil {
			return err
		}

		// возвращаем резервы на склад
		for _, reserveItem := range needToDropReserve.ReserveItems {
			err := m.repository.lomsRepository.ReturnReserve(ctx, reserveItem)
			if err != nil {
				return err
			}
		}

		// очищаем резервы
		err = m.repository.lomsRepository.ClearReserves(ctx, req.OrderId)
		if err != nil {
			return err
		}

		// вызываем метод смены статуса
		// cancelled
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
