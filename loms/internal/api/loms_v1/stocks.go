package loms_v1

import (
	"context"
	desc "route256/loms/pkg/loms_v1"
)

func (i *Implementation) Stocks(ctx context.Context, req *desc.StocksRequest) (*desc.StocksResponse, error) {
	return &desc.StocksResponse{
		Stocks: []*desc.StockItem{
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
