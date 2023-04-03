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

func (r *repository) MessagesToSend(ctx context.Context) ([]*model.OrderMessage, error) {
	// получаем сообщения, которые ещё не отправлены
	logger.Debug("loms repository", zap.String("handler", "MessagesToSend"))

	db := r.queryEngineProvider.GetQueryEngine(ctx)

	pgBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// по каждому orders_id нам нужно только самое раннее не отправленное сообщение
	query := pgBuilder.Select("distinct on(orders_id) orders_id", "id", "payload", "created_at").
		From(outboxOrdersTable).
		Where("sent = false").
		OrderBy("orders_id", "created_at ASC")

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var messages []*schema.OrderMessage
	if err := pgxscan.Select(ctx, db, &messages, rawQuery, args...); err != nil {
		return nil, err
	}

	return converter.FromRepositoryToMolelOrderMessageSlice(messages), nil
}
