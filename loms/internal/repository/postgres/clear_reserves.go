package repository

import (
	"context"
	"route256/libs/logger"
	"route256/libs/metrics"

	sq "github.com/Masterminds/squirrel"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

func (r *repository) ClearReserves(ctx context.Context, orderId int64) error {
	// убираем товары из резерва
	logger.Debug("loms repository", zap.String("handler", "ClearReserves"), zap.Int64("orderId", orderId))

	span, ctx := opentracing.StartSpanFromContext(ctx, "loms repository ClearReserves processing")
	defer span.Finish()

	span.SetTag("order_id", orderId)

	db := r.queryEngineProvider.GetQueryEngine(ctx)

	pgBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	metrics.QueryCounter.WithLabelValues("delete", itemsStocksReservationTable).Inc()

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
