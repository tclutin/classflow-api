package auth

import (
	"context"
	"fmt"
	"github.com/tclutin/classflow-api/internal/config"
	"github.com/tclutin/classflow-api/internal/domain/errors"
	domenErr "github.com/tclutin/classflow-api/internal/domain/errors"
	"github.com/tclutin/classflow-api/internal/domain/user"
	"github.com/tclutin/classflow-api/pkg/hash"
	"github.com/tclutin/classflow-api/pkg/jwt"
	"time"
)

type UserService interface {
	GetById(ctx context.Context, userID uint64) (user.User, error)
	GetByEmail(ctx context.Context, email string) (user.User, error)
	GetByTelegramChatId(ctx context.Context, telegramChatID int64) (user.User, error)
	Create(ctx context.Context, user user.User) (uint64, error)
}

type Service struct {
	userService  UserService
	tokenManager jwt.Manager
	cfg          *config.Config
}

func NewService(
	userService UserService,
	tokenManager jwt.Manager,
	cfg *config.Config,
) *Service {

	return &Service{
		userService:  userService,
		tokenManager: tokenManager,
		cfg:          cfg,
	}
}

func (s *Service) SignUp(ctx context.Context, dto SignUpDTO) (TokenDTO, error) {
	_, err := s.userService.GetByEmail(ctx, dto.Email)
	if err == nil {
		return TokenDTO{}, domenErr.ErrUserAlreadyExists
	}

	bcryptHash, err := hash.NewBcryptHash(dto.Password)
	if err != nil {
		return TokenDTO{}, fmt.Errorf("failed to hash password: %w", err)
	}

	entity := user.User{
		Email:                &dto.Email,
		PasswordHash:         &bcryptHash,
		Role:                 user.Admin,
		FullName:             nil,
		TelegramChatID:       nil,
		TelegramUsername:     nil,
		NotificationDelay:    nil,
		NotificationsEnabled: nil,
		CreatedAt:            time.Now(),
	}

	userID, err := s.userService.Create(ctx, entity)
	if err != nil {
		return TokenDTO{}, err
	}

	token, err := s.tokenManager.NewToken(userID, s.cfg.JWT.Expire)
	if err != nil {
		return TokenDTO{}, fmt.Errorf("failed to create access token: %w", err)
	}

	return TokenDTO{
		AccessToken:  token,
		RefreshToken: "another time...",
	}, nil
}

func (s *Service) LogIn(ctx context.Context, dto LogInDTO) (TokenDTO, error) {
	usr, err := s.userService.GetByEmail(ctx, dto.Email)
	if err != nil {
		return TokenDTO{}, err
	}

	if !hash.CompareBcryptHash(*usr.PasswordHash, dto.Password) {
		return TokenDTO{}, errors.ErrWrongPassword
	}

	token, err := s.tokenManager.NewToken(usr.UserID, s.cfg.JWT.Expire)
	if err != nil {
		return TokenDTO{}, fmt.Errorf("failed to create access token: %w", err)
	}

	return TokenDTO{
		AccessToken:  token,
		RefreshToken: "another time...",
	}, nil
}

func (s *Service) SignUpWithTelegram(ctx context.Context, dto SignUpWithTelegramDTO) (TokenDTO, error) {
	_, err := s.userService.GetByTelegramChatId(ctx, dto.TelegramChatID)
	if err == nil {
		return TokenDTO{}, errors.ErrUserAlreadyExists
	}

	notificationsEnabled := false

	entity := user.User{
		Email:                nil,
		PasswordHash:         nil,
		Role:                 user.Student,
		FullName:             &dto.Fullname,
		TelegramUsername:     &dto.TelegramUsername,
		TelegramChatID:       &dto.TelegramChatID,
		NotificationDelay:    nil,
		NotificationsEnabled: &notificationsEnabled,
		CreatedAt:            time.Now(),
	}

	userID, err := s.userService.Create(ctx, entity)
	if err != nil {
		return TokenDTO{}, err
	}

	token, err := s.tokenManager.NewToken(userID, s.cfg.JWT.Expire)
	if err != nil {
		return TokenDTO{}, fmt.Errorf("failed to create access token: %w", err)
	}

	return TokenDTO{
		AccessToken:  token,
		RefreshToken: "another time...",
	}, nil
}

func (s *Service) LogInWithTelegramRequest(ctx context.Context, dto LogInWithTelegramDTO) (TokenDTO, error) {
	usr, err := s.userService.GetByTelegramChatId(ctx, dto.TelegramChatID)
	if err != nil {
		return TokenDTO{}, err
	}

	token, err := s.tokenManager.NewToken(usr.UserID, s.cfg.JWT.Expire)
	if err != nil {
		return TokenDTO{}, fmt.Errorf("failed to create access token: %w", err)
	}

	return TokenDTO{
		AccessToken:  token,
		RefreshToken: "another time...",
	}, nil
}

func (s *Service) Who(ctx context.Context, userID uint64) (user.User, error) {
	usr, err := s.userService.GetById(ctx, userID)
	if err != nil {
		return usr, err
	}

	return usr, nil
}

func (s *Service) VerifyAndGetCredentials(ctx context.Context, token string) (user.User, error) {
	var user user.User

	userID, err := s.tokenManager.ParseToken(token)
	if err != nil {
		return user, err
	}

	user, err = s.userService.GetById(ctx, userID)
	if err != nil {
		return user, err
	}

	return user, nil
}
