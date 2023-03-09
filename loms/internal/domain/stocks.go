package domain

import (
	"context"
	"route256/loms/internal/domain/model"
)

func (s *service) Stocks(ctx context.Context, req *model.StocksRequest) (*model.StocksResponse, error) {
	// получаем остатки на складах
	response, err := s.repository.lomsRepository.Stocks(ctx, req)
	if err != nil {
		return nil, err
	}

	// нужна проверка, если ничего не вернулось

	return response, nil
}
