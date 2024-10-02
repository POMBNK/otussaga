package order

import (
	"context"
	"github.com/POMBNK/orderservice/internal/entity"
)

type OrderRepo interface {
	InsertOrder(ctx context.Context, order entity.Order) (int, error)
	RollbackOrder(ctx context.Context, orderID int) error
}
