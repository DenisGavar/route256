package domain

import (
	"context"
	"route256/loms/internal/domain/model"

	"github.com/pkg/errors"
)

func (s *service) ListOrder(ctx context.Context, req *model.ListOrderRequest) (*model.ListOrderResponse, error) {
	// получаем список товаров заказа

	response, err := s.repository.lomsRepository.ListOrder(ctx, req)
	if err != nil {
		return nil, errors.WithMessage(err, "getting list order")
	}

	return response, nil
}
