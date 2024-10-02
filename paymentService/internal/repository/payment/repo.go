package payment

import (
	"context"
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/POMBNK/paymentService/internal/entity"
	"github.com/POMBNK/paymentService/pkg/client/postgres"
)

type Repo struct {
	conn postgres.Client
}

func NewRepo(conn postgres.Client) *Repo {
	return &Repo{conn: conn}
}

func (r *Repo) CreatePayment(ctx context.Context, payment entity.Payment) (int, error) {
	insertBuilder := sq.Insert("payments").PlaceholderFormat(sq.Dollar).
		Columns("order_id").
		Values(payment.OrderID).
		Suffix("RETURNING id")

	query, args, err := insertBuilder.ToSql()
	if err != nil {
		return 0, err
	}

	var id int
	err = r.conn.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repo) IsPaymentAlreadyExist(ctx context.Context, payment entity.Payment) (bool, error) {
	selectBuilder := sq.Select("order_id").PlaceholderFormat(sq.Dollar).
		From("payments").
		Where(sq.Eq{"order_id": payment.OrderID})

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return true, err
	}

	var orderID int
	err = r.conn.QueryRow(ctx, query, args...).Scan(&orderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return true, err
	}

	if orderID != 0 {
		return true, nil
	}

	return false, nil
}

func (r *Repo) RollbackPayment(ctx context.Context, orderID int) error {

	deleteBuilder := sq.Update("payments").PlaceholderFormat(sq.Dollar).
		Set("status", "FAILED").
		Where(sq.Eq{"order_id": orderID})

	query, args, err := deleteBuilder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.conn.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
