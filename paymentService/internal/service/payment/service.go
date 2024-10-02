package payment

import (
	"context"
	"github.com/POMBNK/paymentService/internal/entity"
)

type Service struct {
	EventSender *Sender
	repo        PayRepo
}

func NewService(repo PayRepo, eventSender *Sender) *Service {
	return &Service{repo: repo, EventSender: eventSender}
}

func (s *Service) CreatePayment(ctx context.Context, payment entity.Payment) (int, error) {
	return s.repo.CreatePayment(ctx, payment)
}

func (s *Service) IsPaymentAlreadyExist(ctx context.Context, payment entity.Payment) (bool, error) {
	return s.repo.IsPaymentAlreadyExist(ctx, payment)
}

func (s *Service) RollbackPayment(ctx context.Context, orderID int) error {
	return s.repo.RollbackPayment(ctx, orderID)
}
