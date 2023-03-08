package repository

import (
	"context"
	"route256/loms/internal/converter"
	"route256/loms/internal/domain/model"
	"route256/loms/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
)

func (r *repo) Stocks(ctx context.Context, stocksRequest *model.StocksRequest) (*model.StocksResponse, error) {
	db := r.queryEngineProvider.GetQueryEngine(ctx)

	pgBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := pgBuilder.Select("warehouse_id", "count").
		From(itemsStocksTable).
		Where("sku = ?", stocksRequest.Sku)

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var stockItems []*schema.StockItem
	if err := pgxscan.Select(ctx, db, &stockItems, rawQuery, args...); err != nil {
		return nil, err
	}

	return converter.FromRepositoryToMolelStocksResponse(stockItems), nil
}
