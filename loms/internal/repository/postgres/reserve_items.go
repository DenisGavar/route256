package repository

import (
	"context"
	"log"
	"route256/loms/internal/domain/model"
	"route256/loms/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *repo) ReserveItems(ctx context.Context, orderId int64, req *model.ReserveStocksItem) error {
	log.Printf("reserve: %+v", req)

	// резервируем товары на складах

	db := r.queryEngineProvider.GetQueryEngine(ctx)

	pgBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// получаем строки из БД, которые относятся к данному запросу
	query := pgBuilder.Select("id", "count").
		From(itemsStocksTable).
		Where("warehouse_id = ?", req.WarehouseId).
		Where("sku = ?", req.Sku)

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	var stockItems []*schema.StockItem
	if err := pgxscan.Select(ctx, db, &stockItems, rawQuery, args...); err != nil {
		return err
	}

	// удаляем, пока есть возможность и пока нужно
	countToDelete := req.Count
	for _, stockItem := range stockItems {
		// проверяем, надо ли ещё удалять
		if countToDelete != 0 {
			// надо удалить строку или обновить
			if countToDelete < stockItem.Count {
				// надо убрать только часть количества
				query := pgBuilder.Update(itemsStocksTable).
					Set("count", sq.Expr("count - ?", countToDelete)).
					Where("id = ?", stockItem.StockId)

				rawQuery, args, err := query.ToSql()
				if err != nil {
					return err
				}

				_, err = db.Exec(ctx, rawQuery, args...)
				if err != nil {
					return err
				}
				countToDelete = 0
			} else {
				// надо убрать всё количество
				query := pgBuilder.Delete(itemsStocksTable).
					Where("id = ?", stockItem.StockId)

				rawQuery, args, err := query.ToSql()
				if err != nil {
					return err
				}

				_, err = db.Exec(ctx, rawQuery, args...)
				if err != nil {
					return err
				}
				countToDelete -= stockItem.Count
			}
		} else {
			// всё удалили
			break
		}
	}

	// добавляем запись в резерв
	queryInsert := pgBuilder.Insert(itemsStocksReservationTable).
		Columns("sku", "warehouse_id", "orders_id", "count").
		Values(req.Sku, req.WarehouseId, orderId, req.Count)

	rawQuery, args, err = queryInsert.ToSql()
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, rawQuery, args...)
	if err != nil {
		return err
	}

	return nil
}
