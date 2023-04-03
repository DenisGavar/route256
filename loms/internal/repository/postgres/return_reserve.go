package repository

import (
	"context"
	"fmt"
	"route256/libs/logger"
	"route256/loms/internal/domain/model"

	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

func (r *repository) ReturnReserve(ctx context.Context, reserveStocksItem *model.ReserveStocksItem) error {
	// возвращаем резервы на склад
	logger.Debug("loms repository", zap.String("handler", "ReturnReserve"), zap.String("request", fmt.Sprintf("%+v", reserveStocksItem)))

	db := r.queryEngineProvider.GetQueryEngine(ctx)

	pgBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := pgBuilder.Insert(itemsStocksTable).
		Columns("sku", "warehouse_id", "count").
		Values(reserveStocksItem.Sku, reserveStocksItem.WarehouseId, reserveStocksItem.Count)

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
