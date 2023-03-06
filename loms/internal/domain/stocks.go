package domain

import (
	"context"
	"route256/loms/internal/domain/model"
)

func (m *service) Stocks(ctx context.Context, req *model.StocksRequest) (*model.StocksResponse, error) {
	return &model.StocksResponse{
		Stocks: []*model.StockItem{
			{
				WarehouseId: 5,
				Count:       4,
			},
			{
				WarehouseId: 6,
				Count:       2,
			},
		},
	}, nil
}
