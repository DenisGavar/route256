package repository

import (
	"context"
	"fmt"
	"route256/checkout/internal/domain/model"
	"route256/libs/logger"
	"route256/libs/metrics"

	sq "github.com/Masterminds/squirrel"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

func (r *repository) AddToCart(ctx context.Context, addToCartRequest *model.AddToCartRequest) error {
	// добавляем товары в корзину
	logger.Debug("checkout repository", zap.String("handler", "AddToCart"), zap.String("request", fmt.Sprintf("%+v", addToCartRequest)))

	span, ctx := opentracing.StartSpanFromContext(ctx, "checkout repository AddToCart processing")
	defer span.Finish()

	span.SetTag("user", addToCartRequest.User)
	span.SetTag("sku", addToCartRequest.Sku)

	db := r.queryEngineProvider.GetQueryEngine(ctx)

	pgBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	metrics.QueryCounter.WithLabelValues("insert", basketsTable).Inc()

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
