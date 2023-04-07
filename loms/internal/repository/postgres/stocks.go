package repository

import (
	"context"
	"fmt"
	"route256/libs/logger"
	"route256/libs/metrics"
	"route256/loms/internal/converter"
	"route256/loms/internal/domain/model"
	"route256/loms/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

func (r *repository) Stocks(ctx context.Context, stocksRequest *model.StocksRequest) (*model.StocksResponse, error) {
	// получаем товары на складах
	logger.Debug("loms repository", zap.String("handler", "Stocks"), zap.String("request", fmt.Sprintf("%+v", stocksRequest)))

	span, ctx := opentracing.StartSpanFromContext(ctx, "loms repository Stocks processing")
	defer span.Finish()

	span.SetTag("sku", stocksRequest.Sku)

	db := r.queryEngineProvider.GetQueryEngine(ctx)

	pgBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	metrics.QueryCounter.WithLabelValues("select", itemsStocksTable).Inc()

	query := pgBuilder.Select("warehouse_id", "sum(count) count").
		From(itemsStocksTable).
		Where("sku = ?", stocksRequest.Sku).
		GroupBy("warehouse_id")

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
