package domain

import (
	"context"
	"route256/loms/internal/domain/model"

	"github.com/pkg/errors"
)

func (s *service) MessagesToSend(ctx context.Context) ([]*model.OrderMessage, error) {
	// получаем сообщения для отправки

	response, err := s.repository.lomsRepository.MessagesToSend(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, ErrGettingMessagesToSend.Error())
	}

	return response, nil
}
