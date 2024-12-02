package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	domenErr "github.com/tclutin/classflow-api/internal/domain/errors"
)

type Repository interface {
	Create(ctx context.Context, user User) (uint64, error)
	Update(ctx context.Context, user User) error
	GetById(ctx context.Context, userID uint64) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	GetByTelegramChatId(ctx context.Context, telegramChatID int64) (User, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, user User) (uint64, error) {
	userID, err := s.repo.Create(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return userID, nil
}

func (s *Service) Update(ctx context.Context, user User) error {
	return s.repo.Update(ctx, user)
}

func (s *Service) UpdatePartial(ctx context.Context, dto PartialUpdateUserDTO, userID uint64) error {
	user, err := s.GetById(ctx, userID)
	if err != nil {
		return err
	}

	if dto.FullName != nil {
		user.FullName = dto.FullName
	}

	if dto.NotificationDelay != nil {
		user.NotificationDelay = dto.NotificationDelay
	}

	if dto.NotificationsEnabled != nil {
		user.NotificationsEnabled = dto.NotificationsEnabled
	}

	return s.Update(ctx, user)
}

func (s *Service) GetByTelegramChatId(ctx context.Context, telegramChatID int64) (User, error) {
	user, err := s.repo.GetByTelegramChatId(ctx, telegramChatID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, domenErr.ErrUserNotFound
		}
		return user, fmt.Errorf("failed to get user by tgchatid: %w", err)
	}

	return user, nil
}

func (s *Service) GetByEmail(ctx context.Context, email string) (User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, domenErr.ErrUserNotFound
		}

		return user, fmt.Errorf("failed to get user by email: %w", err)
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
