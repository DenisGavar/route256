package domain

import (
	"context"
	"fmt"
	"route256/libs/logger"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (s *service) MessageSent(ctx context.Context, id int64) error {
	// помечаем сообщение отправленным
	logger.Debug("loms domain", zap.String("handler", "MessageSent"), zap.String("id", fmt.Sprintf("%+v", id)))

	span, ctx := opentracing.StartSpanFromContext(ctx, "loms domain MessageSent processing")
	defer span.Finish()

	span.SetTag("outbox_orders_id", id)

	err := s.repository.lomsRepository.MessageSent(ctx, id)
	if err != nil {
		return errors.WithMessage(err, ErrChangingMessageSent.Error())
	}

	return nil
}
