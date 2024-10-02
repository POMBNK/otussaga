package order

import (
	"context"
	"github.com/POMBNK/orderservice/internal/entity"
	"github.com/POMBNK/orderservice/internal/repository/tx"
	"github.com/POMBNK/orderservice/internal/service/event"
)

type Service struct {
	repo         OrderRepo
	eventService *event.Service
	eventSender  *event.Sender

	executer tx.Txer
}

func NewService(repo OrderRepo, eventService *event.Service, eventSender *event.Sender, executer tx.Txer) *Service {
	return &Service{repo: repo, eventService: eventService, eventSender: eventSender, executer: executer}
}

func (s *Service) CreateOrder(ctx context.Context, order entity.Order) (int, error) {

	var id int
	err := s.executer.WithTx(ctx, func(ctx context.Context) error {
		var txErr error

		id, txErr = s.repo.InsertOrder(ctx, order)
		if txErr != nil {
			return txErr
		}

		order.ID = id
		orderEvent, err := entity.NewEvent[entity.Order](entity.EventInfo{Type: entity.NewOrderEventType}, order)
		if err != nil {
			return err
		}

		orderEvent.Info.ID = id
		id, txErr = s.eventService.CreateEvent(ctx, orderEvent)
		if txErr != nil {
			return txErr
		}

		orderEvent.Info.ID = id
		txErr = s.eventSender.SendNewEvent(ctx, orderEvent)
		if txErr != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Service) RollbackOrder(ctx context.Context, orderID int) error {
	return s.repo.RollbackOrder(ctx, orderID)
}
