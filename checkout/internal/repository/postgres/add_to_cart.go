package repository

import (
	"context"
	"fmt"
	"route256/checkout/internal/domain/model"

	sq "github.com/Masterminds/squirrel"
)

func (r *repo) AddToCart(ctx context.Context, addToCartRequest *model.AddToCartRequest) error {
	// добавляем товары в корзину
	db := r.queryEngineProvider.GetQueryEngine(ctx)

	pgBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := pgBuilder.Insert(basketsTable).
		Columns("user_id", "sku", "count").
		Values(addToCartRequest.User, addToCartRequest.Sku, addToCartRequest.Count).
		Suffix(fmt.Sprintf("ON CONFLICT (user_id,sku) DO UPDATE SET count = %s.count + ?", basketsTable), addToCartRequest.Count)

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