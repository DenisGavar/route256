package domain

import (
	"context"
	"route256/loms/internal/domain/model"
)

func (m *service) Stocks(ctx context.Context, req *model.StocksRequest) (*model.StocksResponse, error) {
	response, err := m.repository.lomsRepository.Stocks(ctx, req)
	if err != nil {
		return nil, err
	}

	// нужна проверка, если ничего не вернулось

	return response, nil
}
