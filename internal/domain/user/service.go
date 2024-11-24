package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	domenErr "github.com/tclutin/classflow-api/internal/domain/errors"
)

type Repository interface {
	GetById(ctx context.Context, userID uint64) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	Create(ctx context.Context, user User) (uint64, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetByEmail(ctx context.Context, email string) (User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, domenErr.ErrUserNotFound
		}

		return User{}, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}

func (s *Service) GetById(ctx context.Context, userID uint64) (User, error) {
	user, err := s.repo.GetById(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, domenErr.ErrUserNotFound
		}

		return user, fmt.Errorf("failed to get user by id: %w", err)
	}
	return user, nil
}

func (s *Service) Create(ctx context.Context, user User) (uint64, error) {
	_, err := s.GetByEmail(ctx, user.Email)
	if err == nil {
		return 0, domenErr.ErrUserAlreadyExists
	}

	userID, err := s.repo.Create(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return userID, nil
}
