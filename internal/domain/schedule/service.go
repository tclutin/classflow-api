package schedule

import "context"

type Repository interface {
	Create(ctx context.Context, schedule []Schedule) error
	GetSchedulesByGroupId(ctx context.Context, filter FilterDTO, groupID uint64) ([]DetailsScheduleDTO, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, schedule []Schedule) error {
	return s.repo.Create(ctx, schedule)
}

func (s *Service) GetSchedulesByGroupId(ctx context.Context, filter FilterDTO, groupID uint64) ([]DetailsScheduleDTO, error) {
	return s.repo.GetSchedulesByGroupId(ctx, filter, groupID)
}
