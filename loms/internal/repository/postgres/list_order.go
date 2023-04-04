package repository

import (
	"context"
	"fmt"
	"route256/libs/logger"
	"route256/loms/internal/converter"
	"route256/loms/internal/domain/model"
	"route256/loms/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

func (r *repository) ListOrder(ctx context.Context, req *model.ListOrderRequest) (*model.ListOrderResponse, error) {
	// получаем заказ
	logger.Debug("loms repository", zap.String("handler", "ListOrder"), zap.String("request", fmt.Sprintf("%+v", req)))

	span, ctx := opentracing.StartSpanFromContext(ctx, "loms repository ListOrder processing")
	defer span.Finish()

	span.SetTag("order_id", req.OrderId)

	db := r.queryEngineProvider.GetQueryEngine(ctx)

	pgBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// получаем заказ
	query := pgBuilder.Select("user_id", "status").
		From(ordersTable).
		Where("id = ?", req.OrderId)

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var order schema.Order
	if err := pgxscan.Get(ctx, db, &order, rawQuery, args...); err != nil {
		return nil, err
	}

	// получаем сроки заказа
	query = pgBuilder.Select("sku", "sum(count) as count").
		From(orderItemsTable).
		Where("orders_id = ?", req.OrderId).
		GroupBy("sku")

	rawQuery, args, err = query.ToSql()
	if err != nil {
		// возвращать ошибку, что заказ с таким id не найден
		// TODO
		return nil, err
	}

	var orderItems []*schema.OrderItem
	if err := pgxscan.Select(ctx, db, &orderItems, rawQuery, args...); err != nil {
		return nil, err
	}

	return converter.FromRepositoryToMolelListOrderResponse(&order, orderItems), nil

}
