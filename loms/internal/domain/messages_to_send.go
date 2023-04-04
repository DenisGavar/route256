package domain

import (
	"context"
	"route256/libs/logger"
	"route256/loms/internal/domain/model"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (s *service) MessagesToSend(ctx context.Context) ([]*model.OrderMessage, error) {
	// получаем сообщения для отправки
	logger.Debug("loms domain", zap.String("handler", "MessagesToSend"))

	span, ctx := opentracing.StartSpanFromContext(ctx, "loms domain MessagesToSend processing")
	defer span.Finish()

	response, err := s.repository.lomsRepository.MessagesToSend(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, ErrGettingMessagesToSend.Error())
	}

	return response, nil
}
