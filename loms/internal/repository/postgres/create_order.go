package repository

import (
	"context"
	"route256/loms/internal/domain/model"
	"route256/loms/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
)

func (r *repo) CreateOrder(ctx context.Context, req *model.CreateOrderRequest) (*model.CreateOrderResponse, error) {
	db := r.queryEngineProvider.GetQueryEngine(ctx)

	pgBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// создаём заказ
	query := pgBuilder.Insert(ordersTable).
		Columns("user_id", "status").
		Values(req.User, model.OrderStatusNew).
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
