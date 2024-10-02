package postgres

import (
	"errors"
	"github.com/jackc/pgx/v5"
)

var ErrNotFound = errors.New("запись не найдена")

func handleError(err error) error {
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return ErrNotFound
	default:
		return err
	}
}
