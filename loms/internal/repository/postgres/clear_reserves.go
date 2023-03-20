package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

func (r *repo) ClearReserves(ctx context.Context, orderId int64) error {
	// убираем товары из резерва
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
