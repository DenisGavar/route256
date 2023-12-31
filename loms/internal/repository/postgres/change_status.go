package repository

import (
	"context"
	"encoding/json"
	"route256/libs/logger"
	"route256/libs/metrics"
	"route256/loms/internal/domain/model"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

func (r *repository) ChangeStatus(ctx context.Context, orderId int64, status string) error {
	// меняем статус заказа
	logger.Debug("loms repository", zap.String("handler", "ChangeStatus"), zap.Int64("orderId", orderId), zap.String("status", status))

	span, ctx := opentracing.StartSpanFromContext(ctx, "loms repository ChangeStatus processing")
	defer span.Finish()

	span.SetTag("order_id", orderId)
	span.SetTag("status", status)

	changingStatusTime := time.Now()

	db := r.queryEngineProvider.GetQueryEngine(ctx)

	pgBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	metrics.QueryCounter.WithLabelValues("update", ordersTable).Inc()

	query := pgBuilder.Update(ordersTable).
		Set("status", status).
		Set("changed_at", changingStatusTime).
		Where("id = ?", orderId)

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, rawQuery, args...)
	if err != nil {
		return err
	}

	// сохраняем информацию об изменении статуса заказа для дальнейшей отправки в kafka

	// создаём тело сообщения
	req := &model.ListOrderRequest{
		OrderId: orderId,
	}
	res, err := r.ListOrder(ctx, req)
	if err != nil {
		return err
	}
	payload, err := json.Marshal(res)
	if err != nil {
		return err
	}

	metrics.QueryCounter.WithLabelValues("insert", outboxOrdersTable).Inc()

	// сохраняем тело сообщения
	queryInsert := pgBuilder.Insert(outboxOrdersTable).
		Columns("orders_id", "payload", "created_at", "sent").
		Values(orderId, payload, changingStatusTime, false)

	rawQuery, args, err = queryInsert.ToSql()
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, rawQuery, args...)
	if err != nil {
		return err
	}

	return nil
}
