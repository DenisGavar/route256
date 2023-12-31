package transactor

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"

	"go.uber.org/multierr"
)

type QueryEngine interface {
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

type QueryEngineProvider interface {
	GetQueryEngine(ctx context.Context) QueryEngine // tx/pool
}

type TransactionManager interface {
	RunRepeatableRead(ctx context.Context, f func(ctxTX context.Context) error) error
}

type Transactor interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

type DB interface {
	QueryEngine
	Transactor
}

type transactionManager struct {
	pool DB
}

func NewTransactionManager(pool DB) *transactionManager {
	return &transactionManager{
		pool: pool,
	}
}

type key string

const TxKey = key("tx")

func (tm *transactionManager) RunRepeatableRead(ctx context.Context, fx func(ctxTX context.Context) error) error {
	tx, err := tm.pool.BeginTx(ctx,
		pgx.TxOptions{
			IsoLevel: pgx.RepeatableRead,
		})
	if err != nil {
		return err
	}

	if err := fx(context.WithValue(ctx, TxKey, tx)); err != nil {
		return multierr.Combine(err, tx.Rollback(ctx))
	}

	if err := tx.Commit(ctx); err != nil {
		return multierr.Combine(err, tx.Rollback(ctx))
	}

	return nil
}

func (tm *transactionManager) GetQueryEngine(ctx context.Context) QueryEngine {
	tx, ok := ctx.Value(TxKey).(QueryEngine)
	if ok && tx != nil {
		return tx
	}

	return tm.pool
}

func (tm *transactionManager) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return tm.pool.BeginTx(ctx, txOptions)
}
