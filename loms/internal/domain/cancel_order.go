package domain

import (
	"context"
	"fmt"
	"route256/libs/logger"
	"route256/loms/internal/domain/model"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (s *service) CancelOrder(ctx context.Context, req *model.CancelOrderRequest) error {
	// отменяем заказ
	logger.Debug("loms domain", zap.String("handler", "CancelOrder"), zap.String("request", fmt.Sprintf("%+v", req)))

	span, ctx := opentracing.StartSpanFromContext(ctx, "loms domain CancelOrder processing")
	defer span.Finish()

	span.SetTag("order_id", req.OrderId)

	err := s.repository.transactionManager.RunRepeatableRead(ctx, func(ctxTX context.Context) error {
		// получаем резервы, которые надо вернуть на склад
		needToDropReserve, err := s.repository.lomsRepository.Reserves(ctxTX, req.OrderId)
		if err != nil {
			return errors.WithMessage(err, ErrCheckingReserves.Error())
		}

		// возвращаем резервы на склад
		for _, reserveItem := range needToDropReserve.ReserveItems {
			err := s.repository.lomsRepository.ReturnReserve(ctxTX, reserveItem)
			if err != nil {
				return errors.WithMessage(err, ErrReturningReserves.Error())
			}
		}

		// очищаем резервы
		err = s.repository.lomsRepository.ClearReserves(ctxTX, req.OrderId)
		if err != nil {
			return errors.WithMessage(err, ErrClearingReserves.Error())
		}

		return nil
	})

	if err != nil {
		return err
	}

	// вызываем метод смены статуса
	// cancelled
	err = s.repository.lomsRepository.ChangeStatus(ctx, req.OrderId, model.OrderStatusCancelled)
	if err != nil {
		return errors.WithMessage(err, ErrChangingStatus.Error())
	}

	return nil
}
