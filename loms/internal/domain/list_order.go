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

func (s *service) ListOrder(ctx context.Context, req *model.ListOrderRequest) (*model.ListOrderResponse, error) {
	// получаем список товаров заказа
	logger.Debug("loms domain", zap.String("handler", "ListOrder"), zap.String("request", fmt.Sprintf("%+v", req)))

	span, ctx := opentracing.StartSpanFromContext(ctx, "loms domain ListOrder processing")
	defer span.Finish()

	span.SetTag("order_id", req.OrderId)

	response, err := s.repository.lomsRepository.ListOrder(ctx, req)
	if err != nil {
		return nil, errors.WithMessage(err, ErrGettingListOrder.Error())
	}

	return response, nil
}
