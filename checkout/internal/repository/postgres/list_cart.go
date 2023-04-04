package repository

import (
	"context"
	"fmt"
	"route256/checkout/internal/converter"
	"route256/checkout/internal/domain/model"
	"route256/checkout/internal/repository/schema"
	"route256/libs/logger"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

func (r *repository) ListCart(ctx context.Context, listCartRequest *model.ListCartRequest) (*model.ListCartResponse, error) {
	// получаем список товаров в корзине
	logger.Debug("checkout repository", zap.String("handler", "ListCart"), zap.String("request", fmt.Sprintf("%+v", listCartRequest)))

	span, ctx := opentracing.StartSpanFromContext(ctx, "checkout repository ListCart processing")
	defer span.Finish()

	span.SetTag("user", listCartRequest.User)

	db := r.queryEngineProvider.GetQueryEngine(ctx)

	pgBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := pgBuilder.Select("sku", "sum(count) as count").
		From(basketsTable).
		Where("user_id = ?", listCartRequest.User).
		OrderBy("sku").
		GroupBy("sku")

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
