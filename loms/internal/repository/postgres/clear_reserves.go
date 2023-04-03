package repository

import (
	"context"
	"route256/libs/logger"

	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

func (r *repository) ClearReserves(ctx context.Context, orderId int64) error {
	// убираем товары из резерва
	logger.Debug("loms repository", zap.String("handler", "ClearReserves"), zap.Int64("orderId", orderId))

	db := r.queryEngineProvider.GetQueryEngine(ctx)

	pgBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := pgBuilder.Delete(itemsStocksReservationTable).
		Where("orders_id = ?", orderId)

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, rawQuery, args...)
	if err != nil {
		return err
	}

	return nil
}
