package repository

import (
	"context"
	"fmt"
	"route256/libs/logger"
	"route256/libs/metrics"
	"route256/loms/internal/domain/model"
	"route256/loms/internal/repository/schema"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

func (r *repository) CreateOrder(ctx context.Context, req *model.CreateOrderRequest) (*model.CreateOrderResponse, error) {
	// создаём заказ
	logger.Debug("loms repository", zap.String("handler", "CreateOrder"), zap.String("request", fmt.Sprintf("%+v", req)))

	span, ctx := opentracing.StartSpanFromContext(ctx, "loms repository CreateOrder processing")
	defer span.Finish()

	span.SetTag("user", req.User)

	db := r.queryEngineProvider.GetQueryEngine(ctx)

	pgBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	metrics.QueryCounter.WithLabelValues("insert", ordersTable).Inc()

	// создаём заказ
	query := pgBuilder.Insert(ordersTable).
		Columns("user_id", "status", "created_at", "changed_at").
		Values(req.User, model.OrderStatusNew, time.Now(), time.Now()).
		Suffix("RETURNING id")

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	raws, err := db.Query(ctx, rawQuery, args...)
	if err != nil {
		return nil, err
	}
	defer raws.Close()

	var order schema.Order
	if err := pgxscan.ScanOne(&order, raws); err != nil {
		return nil, err
	}

	metrics.QueryCounter.WithLabelValues("insert", orderItemsTable).Inc()

	// создаём строки заказа
	query = pgBuilder.Insert(orderItemsTable).
		Columns("orders_id", "sku", "count")

	for _, orderItem := range req.Items {
		query = query.Values(order.OrderId, orderItem.Sku, orderItem.Count)
	}

	rawQuery, args, err = query.ToSql()
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(ctx, rawQuery, args...)
	if err != nil {
		return nil, err
	}

	return &model.CreateOrderResponse{
		OrderId: order.OrderId,
	}, nil
}
