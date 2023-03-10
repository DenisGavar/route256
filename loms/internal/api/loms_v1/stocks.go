package loms_v1

import (
	"context"
	"route256/loms/internal/converter"
	desc "route256/loms/pkg/loms_v1"
)

func (i *Implementation) Stocks(ctx context.Context, req *desc.StocksRequest) (*desc.StocksResponse, error) {
	//log.Printf("stocks: %+v", req)

	response, err := i.lomsService.Stocks(ctx, converter.FromDescToMolelStocksRequest(req))
	if err != nil {
		return nil, err
	}

	return converter.FromModelToDescStocksResponse(response), nil
}
