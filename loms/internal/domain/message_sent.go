package domain

import (
	"context"

	"github.com/pkg/errors"
)

func (s *service) MessageSent(ctx context.Context, id int64) error {
	// получаем сообщения для отправки

	err := s.repository.lomsRepository.MessageSent(ctx, id)
	if err != nil {
		return errors.WithMessage(err, ErrChangingMessageSent.Error())
	}

	return nil
}
