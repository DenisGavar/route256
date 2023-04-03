package repository

import (
	"context"
	"fmt"
	"route256/checkout/internal/domain/model"
	"route256/libs/logger"

	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

func (r *repository) AddToCart(ctx context.Context, addToCartRequest *model.AddToCartRequest) error {
	// добавляем товары в корзину
	logger.Debug("checkout repository", zap.String("handler", "AddToCart"), zap.String("request", fmt.Sprintf("%+v", addToCartRequest)))

	db := r.queryEngineProvider.GetQueryEngine(ctx)

	pgBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := pgBuilder.Insert(basketsTable).
		Columns("user_id", "sku", "count").
		Values(addToCartRequest.User, addToCartRequest.Sku, addToCartRequest.Count)

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
