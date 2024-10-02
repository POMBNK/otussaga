package event

import (
	"context"
	"github.com/POMBNK/orderservice/internal/entity"
)

type Service struct {
	repo EventRepo
}

func NewService(repo EventRepo) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateEvent(ctx context.Context, event entity.Event) (int, error) {
	return s.repo.InsertEvent(ctx, event)
}
