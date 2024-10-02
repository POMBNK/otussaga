package order

import (
	"context"
	"encoding/json"
	sq "github.com/Masterminds/squirrel"
	"github.com/POMBNK/orderservice/internal/entity"
	"github.com/POMBNK/orderservice/pkg/client/postgres"
)

type Repo struct {
	conn postgres.Client
}

func NewRepo(conn postgres.Client) *Repo {
	return &Repo{conn: conn}
}

func (r *Repo) InsertOrder(ctx context.Context, order entity.Order) (int, error) {

	goodsb, err := json.Marshal(order)
	if err != nil {
		return 0, err
	}

	insertBuilder := sq.Insert("orders").PlaceholderFormat(sq.Dollar).
		Columns("goods").
		Values(goodsb).
		Suffix("RETURNING id")

	query, args, err := insertBuilder.ToSql()
	if err != nil {
		return 0, err
	}

	var orderID int
	err = r.conn.QueryRow(ctx, query, args...).Scan(&orderID)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}

func (r *Repo) RollbackOrder(ctx context.Context, orderID int) error {
	updateBuilder := sq.Update("orders").PlaceholderFormat(sq.Dollar).
		Set("status", "FAILED").
		Where(sq.Eq{"id": orderID})

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.conn.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
