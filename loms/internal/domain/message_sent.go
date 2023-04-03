package domain

import (
	"context"
	"fmt"
	"route256/libs/logger"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (s *service) MessageSent(ctx context.Context, id int64) error {
	// помечаем сообщение отправленным
	logger.Debug("loms domain", zap.String("handler", "MessageSent"), zap.String("id", fmt.Sprintf("%+v", id)))

	err := s.repository.lomsRepository.MessageSent(ctx, id)
	if err != nil {
		return errors.WithMessage(err, ErrChangingMessageSent.Error())
	}

	return nil
}
