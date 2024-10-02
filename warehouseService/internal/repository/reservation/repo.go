package reservation

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/POMBNK/warehouseService/internal/entity"
	"github.com/POMBNK/warehouseService/pkg/client/postgres"
)

type Repo struct {
	conn postgres.Client
}

func NewRepo(conn postgres.Client) *Repo {
	return &Repo{conn: conn}
}

func (r *Repo) InsertReservation(ctx context.Context, reservation entity.Reservation) (int, error) {

	goods, err := json.Marshal(reservation.Goods)
	if err != nil {
		return 0, err
	}

	insertBuilder := sq.Insert("order_reservations").PlaceholderFormat(sq.Dollar).
		Columns("order_id", "goods").
		Values(reservation.OrderID, goods).
		Suffix("RETURNING id")

	query, args, err := insertBuilder.ToSql()
	if err != nil {
		return 0, err
	}

	var reservationID int
	err = r.conn.QueryRow(ctx, query, args...).Scan(&reservationID)
	if err != nil {
		return 0, err
	}

	return reservationID, nil
}

func (r *Repo) IsReservationAlreadyExist(ctx context.Context, reservation entity.Reservation) (bool, error) {
	selectBuilder := sq.Select("order_id").PlaceholderFormat(sq.Dollar).
		From("order_reservations").
		Where(sq.Eq{"order_id": reservation.OrderID})

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

func (r *Repo) RollbackReservation(ctx context.Context, orderID int) error {
	updateBuilder := sq.Update("order_reservations").PlaceholderFormat(sq.Dollar).
		Set("status", "FAILED").
		Where(sq.Eq{"order_id": orderID})

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
