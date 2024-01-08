package domain

import (
	"context"
	"fmt"
	"route256/checkout/internal/domain/model"
	"route256/libs/logger"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (s *service) DeleteFromCart(ctx context.Context, req *model.DeleteFromCartRequest) error {
	// удаляем товар из корзины
	logger.Debug("checkout domain", zap.String("handler", "DeleteFromCart"), zap.String("request", fmt.Sprintf("%+v", req)))

	span, ctx := opentracing.StartSpanFromContext(ctx, "checkout domain DeleteFromCart processing")
	defer span.Finish()

	span.SetTag("user", req.User)
	span.SetTag("sku", req.Sku)

	err := s.repository.transactionManager.RunRepeatableRead(ctx, func(ctxTX context.Context) error {
		// удаление товара в нужном количестве из корзины
		err := s.repository.checkoutRepository.DeleteFromCart(ctxTX, req)
		if err != nil {
			return errors.WithMessage(err, ErrDeletingFromCart.Error())
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
