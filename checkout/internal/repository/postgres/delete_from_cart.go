package repository

import (
	"context"
	"route256/checkout/internal/domain/model"

	sq "github.com/Masterminds/squirrel"
)

func (r *repo) DeleteFromCart(ctx context.Context, part bool, deleteFromRequest *model.DeleteFromCartRequest) error {
	// убираем товары из корзины
	db := r.queryEngineProvider.GetQueryEngine(ctx)

	pgBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	if part {
		query := pgBuilder.Update(basketsTable).
			Set("count", sq.Expr("count - ?", deleteFromRequest.Count)).
			Where("user_id = ?", deleteFromRequest.User).
			Where("sku = ?", deleteFromRequest.Sku)

		rawQuery, args, err := query.ToSql()
		if err != nil {
			return err
		}

		_, err = db.Exec(ctx, rawQuery, args...)
		if err != nil {
			return err
		}
	} else {
		query := pgBuilder.Delete(basketsTable).
			Where("user_id = ?", deleteFromRequest.User).
			Where("sku = ?", deleteFromRequest.Sku)

		rawQuery, args, err := query.ToSql()
		if err != nil {
			return err
		}

		_, err = db.Exec(ctx, rawQuery, args...)
		if err != nil {
			return err
		}
	}

	return nil
}
