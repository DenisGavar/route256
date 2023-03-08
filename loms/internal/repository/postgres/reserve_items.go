package repository

import (
	"context"
	"log"
	"route256/loms/internal/domain/model"

	sq "github.com/Masterminds/squirrel"
)

func (r *repo) ReserveItems(ctx context.Context, orderId int64, warehouseId int64, req *model.ReserveStocksItem) error {
	log.Printf("reserve: %+v, %d", req, warehouseId)

	// резервируем товары на складах

	db := r.queryEngineProvider.GetQueryEngine(ctx)

	pgBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	if !req.Part {
		// если part == false, то удаляем запись со склада
		query := pgBuilder.Delete(itemsStocksTable).
			Where("sku = ?", req.Sku).
			Where("warehouse_id = ?", warehouseId)

		rawQuery, args, err := query.ToSql()
		if err != nil {
			return err
		}

		_, err = db.Exec(ctx, rawQuery, args...)
		if err != nil {
			return err
		}
	} else {
		// если part == true, то модифицируем запись на складе
		query := pgBuilder.Update(itemsStocksTable).
			Set("count", sq.Expr("count - ?", req.Count)).
			Where("sku = ?", req.Sku).
			Where("warehouse_id = ?", warehouseId)

		rawQuery, args, err := query.ToSql()
		if err != nil {
			return err
		}

		_, err = db.Exec(ctx, rawQuery, args...)
		if err != nil {
			return err
		}
	}

	// добавляем запись в резерв
	query := pgBuilder.Insert(itemsStocksReservationTable).
		Columns("sku", "warehouse_id", "orders_id", "count").
		Values(req.Sku, warehouseId, orderId, req.Count)

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
