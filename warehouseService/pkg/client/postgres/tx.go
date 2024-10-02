package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

const (
	txKey = "tx"
)

type tx struct {
	Tx       pgx.Tx
	IsActive bool
}

func withTx(ctx context.Context, txx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey, &tx{
		Tx:       txx,
		IsActive: true,
	})
}

func fromContext(ctx context.Context) *tx {
	txx, ok := ctx.Value(txKey).(*tx)
	if !ok {
		return nil
	}
	return txx
}

func (d *Db) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {

	// todo: Добавить в рабочий проект.
	// бага в рабочем проекте оказывается тут :)
	nestedTx := fromContext(ctx)
	if nestedTx != nil && nestedTx.IsActive {
		return fn(ctx)
	}

	txx, err := d.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	tCtx := withTx(ctx, txx)
	defer func() {
		expired := fromContext(tCtx)
		if expired != nil {
			expired.IsActive = false
		}
	}()

	if err := fn(tCtx); err != nil {
		terr := txx.Rollback(tCtx)
		if terr != nil {
			return fmt.Errorf("%w: rollback error: %v", err, terr)
		}
		return err
	}

	return txx.Commit(tCtx)
}
