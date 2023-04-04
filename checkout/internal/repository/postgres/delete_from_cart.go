package repository

import (
	"context"
	"fmt"
	"route256/checkout/internal/domain/model"
	"route256/checkout/internal/repository/schema"
	"route256/libs/logger"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

func (r *repository) DeleteFromCart(ctx context.Context, deleteFromCartRequest *model.DeleteFromCartRequest) error {
	// убираем товары из корзины
	logger.Debug("checkout repository", zap.String("handler", "DeleteFromCart"), zap.String("request", fmt.Sprintf("%+v", deleteFromCartRequest)))

	span, ctx := opentracing.StartSpanFromContext(ctx, "checkout repository DeleteFromCart processing")
	defer span.Finish()

	span.SetTag("user", deleteFromCartRequest.User)
	span.SetTag("sku", deleteFromCartRequest.Sku)

	db := r.queryEngineProvider.GetQueryEngine(ctx)

	pgBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// получаем строки из БД, которые относятся к данному запросу
	query := pgBuilder.Select("id", "count").
		From(basketsTable).
		Where("user_id = ?", deleteFromCartRequest.User).
		Where("sku = ?", deleteFromCartRequest.Sku)

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	var cartItems []*schema.CartItem
	if err := pgxscan.Select(ctx, db, &cartItems, rawQuery, args...); err != nil {
		return err
	}

	// удаляем, пока есть возможность и пока нужно
	countToDelete := deleteFromCartRequest.Count
	for _, cartItem := range cartItems {
		// проверяем, надо ли ещё удалять
		if countToDelete != 0 {
			// надо удалить строку или обновить
			if countToDelete < cartItem.Count {
				// надо убрать только часть количества
				query := pgBuilder.Update(basketsTable).
					Set("count", sq.Expr("count - ?", countToDelete)).
					Where("id = ?", cartItem.BasketId)

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
				query := pgBuilder.Delete(basketsTable).
					Where("id = ?", cartItem.BasketId)

				rawQuery, args, err := query.ToSql()
				if err != nil {
					return err
				}

				_, err = db.Exec(ctx, rawQuery, args...)
				if err != nil {
					return err
				}
				countToDelete -= cartItem.Count
			}
		} else {
			// всё удалили
			break
		}
	}

	return nil
}
