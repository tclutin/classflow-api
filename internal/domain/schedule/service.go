package schedule

import "context"

type Repository interface {
	GetSchedulesByGroupId(ctx context.Context, filter FilterDTO, groupID uint64) ([]DetailsScheduleDTO, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetSchedulesByGroupId(ctx context.Context, filter FilterDTO, groupID uint64) ([]DetailsScheduleDTO, error) {
	return s.repo.GetSchedulesByGroupId(ctx, filter, groupID)
}
