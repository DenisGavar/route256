package repository

import (
	"context"
	"route256/loms/internal/domain/model"
	"route256/loms/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *repo) ReserveItems(ctx context.Context, req *model.CreateOrderRequest) error {

	// резервируем товары на складах, на которых они есть

	db := r.queryEngineProvider.GetQueryEngine(ctx)

	pgBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// получаем где у нас есть товар
	skus := make([]uint32, 0, len(req.Items))
	for _, orderItem := range req.Items {
		skus = append(skus, orderItem.Sku)
	}

	query := pgBuilder.Select("sku", "warehouse_id", "count").
		From(itemsStocksTable).
		Where(sq.Eq{"sku": skus})

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	var stockItems []schema.StockItem
	if err := pgxscan.Select(ctx, db, &stockItems, rawQuery, args...); err != nil {
		return err
	}

	// анализируем сколько откуда надо списать
	// собираем запрос на списание
	// собираем запрос на добавление в резерв

	return nil
}
