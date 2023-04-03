package repository

import (
	"context"
	"route256/libs/logger"

	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

func (r *repository) MessageSent(ctx context.Context, id int64) error {
	// помечаем сообщение отправленным
	logger.Debug("loms repository", zap.String("handler", "MessageSent"), zap.Int64("id", id))

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