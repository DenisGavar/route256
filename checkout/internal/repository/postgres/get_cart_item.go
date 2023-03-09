package repository

import (
	"context"
	"route256/checkout/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *repo) GetCartItemCount(ctx context.Context, userId int64, sku uint32) (uint32, error) {
	// убираем товары из корзины
	db := r.queryEngineProvider.GetQueryEngine(ctx)

	pgBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := pgBuilder.Select("count").
		From(basketsTable).
		Where("user_id = ?", userId).
		Where("sku = ?", sku)

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return 0, err
	}

	raws, err := db.Query(ctx, rawQuery, args...)
	if err != nil {
		return 0, err
	}
	defer raws.Close()

	var order schema.CartItem
	if err := pgxscan.ScanOne(&order, raws); err != nil {
		return 0, err
	}

	return order.Count, nil
}
