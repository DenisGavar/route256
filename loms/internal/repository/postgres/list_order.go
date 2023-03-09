package repository

import (
	"context"
	"route256/loms/internal/converter"
	"route256/loms/internal/domain/model"
	"route256/loms/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *repo) ListOrder(ctx context.Context, req *model.ListOrderRequest) (*model.ListOrderResponse, error) {
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
	query = pgBuilder.Select("sku", "count").
		From(orderItemsTable).
		Where("orders_id = ?", req.OrderId)

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
