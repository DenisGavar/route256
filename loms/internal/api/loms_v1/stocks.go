package loms_v1

import (
	"context"
	"log"
	desc "route256/loms/pkg/loms_v1"
)

func (i *Implementation) Stocks(ctx context.Context, req *desc.StocksRequest) (*desc.StocksResponse, error) {
	log.Printf("stocks: %+v", req)

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
