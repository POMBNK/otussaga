package tx

import (
	"context"
	"github.com/POMBNK/orderservice/pkg/client/postgres"
)

type Txer interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type Tx struct {
	t *postgres.Db
}

func (t *Tx) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return t.t.WithTx(ctx, fn)
}

func New(t *postgres.Db) *Tx {
	return &Tx{t: t}
}
