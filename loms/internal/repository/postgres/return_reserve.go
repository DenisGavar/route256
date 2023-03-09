package repository

import (
	"context"
	"fmt"
	"route256/loms/internal/domain/model"

	sq "github.com/Masterminds/squirrel"
)

func (r *repo) ReturnReserve(ctx context.Context, reserveStocksItem *model.ReserveStocksItem) error {
	// возвращаем резервы на склад
	db := r.queryEngineProvider.GetQueryEngine(ctx)

	pgBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := pgBuilder.Insert(itemsStocksTable).
		Columns("sku", "warehouse_id", "count").
		Values(reserveStocksItem.Sku, reserveStocksItem.WarehouseId, reserveStocksItem.Count).
		Suffix(fmt.Sprintf("ON CONFLICT (sku,warehouse_id) DO UPDATE SET count = %s.count + ?", itemsStocksTable), reserveStocksItem.Count)

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
