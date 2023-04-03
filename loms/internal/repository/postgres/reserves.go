package repository

import (
	"context"
	"route256/libs/logger"
	"route256/loms/internal/converter"
	"route256/loms/internal/domain/model"
	"route256/loms/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"go.uber.org/zap"
)

func (r *repository) Reserves(ctx context.Context, orderId int64) (*model.Reserve, error) {
	// получаем резервы
	logger.Debug("loms repository", zap.String("handler", "Reserves"), zap.Int64("orderId", orderId))

	db := r.queryEngineProvider.GetQueryEngine(ctx)

	pgBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := pgBuilder.Select("sku", "warehouse_id", "sum(count) as count").
		From(itemsStocksReservationTable).
		Where("orders_id = ?", orderId).
		GroupBy("sku", "warehouse_id")

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var reserveItems []*schema.ReserveItem
	if err := pgxscan.Select(ctx, db, &reserveItems, rawQuery, args...); err != nil {
		return nil, err
	}

	return converter.FromRepositoryToMolelReserves(reserveItems), nil
}
