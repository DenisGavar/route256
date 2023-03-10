package domain

import (
	"context"
	"route256/loms/internal/domain/model"
)

func (s *service) ListOrder(ctx context.Context, req *model.ListOrderRequest) (*model.ListOrderResponse, error) {
	response, err := s.repository.lomsRepository.ListOrder(ctx, req)
	if err != nil {
		return nil, err
	}

	return response, nil
}
