package repository

import (
	"context"
	"route256/checkout/internal/converter"
	"route256/checkout/internal/domain/model"
	"route256/checkout/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *repo) ListCart(ctx context.Context, listCartRequest *model.ListCartRequest) (*model.ListCartResponse, error) {
	// получаем список товаров в корзине
	db := r.queryEngineProvider.GetQueryEngine(ctx)

	pgBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := pgBuilder.Select("sku", "count").
		From(basketsTable).
		Where("user_id = ?", listCartRequest.User).
		OrderBy("sku")

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var cartItems []*schema.CartItem
	if err := pgxscan.Select(ctx, db, &cartItems, rawQuery, args...); err != nil {
		return nil, err
	}

	return converter.FromRepositoryToMolelListCartResponse(cartItems), nil
}
