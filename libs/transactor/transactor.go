package transactor

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/multierr"
)

type QueryEngine interface {
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

type QueryEngineProvider interface {
	GetQueryEngine(ctx context.Context) QueryEngine // tx/pool
}

type transactionManager struct {
	pool *pgxpool.Pool
}

func NewTransactionManager(pool *pgxpool.Pool) *transactionManager {
	return &transactionManager{
		pool: pool,
	}
}

type txkey string

const key = txkey("tx")

func (tm *transactionManager) RunRepeatableRead(ctx context.Context, fx func(ctxTX context.Context) error) error {
	tx, err := tm.pool.BeginTx(ctx,
		pgx.TxOptions{
			IsoLevel: pgx.RepeatableRead,
		})
	if err != nil {
		return err
	}

	if err := fx(context.WithValue(ctx, key, tx)); err != nil {
		return multierr.Combine(err, tx.Rollback(ctx))
	}

	if err := tx.Commit(ctx); err != nil {
		return multierr.Combine(err, tx.Rollback(ctx))
	}

	return nil
}

func (tm *transactionManager) GetQueryEngine(ctx context.Context) QueryEngine {
	tx, ok := ctx.Value(key).(QueryEngine)
	if ok && tx != nil {
		return tx
	}

	return tm.pool
}
