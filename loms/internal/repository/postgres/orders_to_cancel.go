package repository

import (
	"context"
	"route256/loms/internal/converter"
	"route256/loms/internal/domain/model"
	"route256/loms/internal/repository/schema"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *repo) OrdersToCancel(ctx context.Context, time time.Time) ([]*model.CancelOrderRequest, error) {
	//получаем заказы на отмену
	db := r.queryEngineProvider.GetQueryEngine(ctx)

	pgBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// получаем заказы, которые:
	// в статусе awaiting_payment
	// время изменения меьше установленного
	query := pgBuilder.Select("id").
		From(ordersTable).
		Where("status = ?", model.OrderStatusAwaitingPayment).
		Where("changed_at < ?", time)

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var orders []*schema.CancelOrderRequest
	if err := pgxscan.Select(ctx, db, &orders, rawQuery, args...); err != nil {
		return nil, err
	}

	return converter.FromRepositoryToMolelCancelOrderRequestSlice(orders), nil
}
