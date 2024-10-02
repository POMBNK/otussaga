package payment

import (
	"context"
	"github.com/POMBNK/paymentService/internal/entity"
)

type PayRepo interface {
	CreatePayment(ctx context.Context, payment entity.Payment) (int, error)
	IsPaymentAlreadyExist(ctx context.Context, payment entity.Payment) (bool, error)
	RollbackPayment(ctx context.Context, orderID int) error
}
