package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

func (r *repository) MessageSent(ctx context.Context, id int64) error {
	// меняем статус сообщения

	db := r.queryEngineProvider.GetQueryEngine(ctx)

	pgBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := pgBuilder.Update(outboxOrdersTable).
		Set("sent", true).
		Where("id = ?", id)

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
