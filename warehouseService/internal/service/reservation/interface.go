package reservation

import (
	"context"
	"github.com/POMBNK/warehouseService/internal/entity"
)

type ReservRepo interface {
	InsertReservation(ctx context.Context, reservation entity.Reservation) (int, error)
	IsReservationAlreadyExist(ctx context.Context, reservation entity.Reservation) (bool, error)
	RollbackReservation(ctx context.Context, orderID int) error
}
