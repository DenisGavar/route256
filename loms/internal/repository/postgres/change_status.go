package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
)

func (r *repo) ChangeStatus(ctx context.Context, orderId int64, status string) error {
	// меняем статус заказа
	db := r.queryEngineProvider.GetQueryEngine(ctx)

	pgBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := pgBuilder.Update(ordersTable).
		Set("status", status).
		Set("changed_at", time.Now()).
		Where("id = ?", orderId)

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
