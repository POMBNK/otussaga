package event

import (
	"context"
	"github.com/POMBNK/orderservice/internal/entity"
)

type EventRepo interface {
	InsertEvent(ctx context.Context, event entity.Event) (int, error)
	GetNewEvents(ctx context.Context) ([]entity.Event, error)
	GetNewEvent(ctx context.Context) (entity.Event, error)
	GetNewEventIDs(ctx context.Context) ([]int, error)
	SetCompleted(ctx context.Context, ids []int) error
}
