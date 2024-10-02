package reservation

import (
	"context"
	"github.com/POMBNK/warehouseService/internal/entity"
	"github.com/POMBNK/warehouseService/internal/repository/tx"
)

type Service struct {
	EventSender *Sender
	executer    tx.Txer
	repo        ReservRepo
}

func NewService(repo ReservRepo, executer tx.Txer, eventSender *Sender) *Service {
	return &Service{repo: repo, executer: executer, EventSender: eventSender}
}

func (s *Service) CreateReservation(ctx context.Context, reservation entity.Reservation) (int, error) {
	var reservationID int
	err := s.executer.WithTx(ctx, func(ctx context.Context) error {
		var txErr error
		reservationID, txErr = s.repo.InsertReservation(ctx, reservation)
		if txErr != nil {
			return txErr
		}

		//txErr = s.EventSender.SendReservationEvent(ctx, reservation)
		//if txErr != nil {
		//	return txErr
		//}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return reservationID, err
}

func (s *Service) IsReservationAlreadyExist(ctx context.Context, reservation entity.Reservation) (bool, error) {
	return s.repo.IsReservationAlreadyExist(ctx, reservation)
}

func (s *Service) RollbackReservation(ctx context.Context, orderID int) error {
	return s.repo.RollbackReservation(ctx, orderID)
}
