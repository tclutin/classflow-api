package domain

import (
	"github.com/tclutin/classflow-api/internal/config"
	"github.com/tclutin/classflow-api/internal/domain/auth"
	"github.com/tclutin/classflow-api/internal/domain/edu"
	"github.com/tclutin/classflow-api/internal/domain/group"
	"github.com/tclutin/classflow-api/internal/domain/schedule"
	"github.com/tclutin/classflow-api/internal/domain/user"
	"github.com/tclutin/classflow-api/internal/repository"
	"github.com/tclutin/classflow-api/pkg/jwt"
	"log/slog"
)

type Services struct {
	Auth     *auth.Service
	User     *user.Service
	Schedule *schedule.Service
	Edu      *edu.Service
	Group    *group.Service
}

func NewServices(
	logger *slog.Logger,
	tokenManager jwt.Manager,
	repositories *repository.Repositories,
	cfg *config.Config,
) *Services {

	userService := user.NewService(repositories.User)
	authService := auth.NewService(userService, tokenManager, cfg)
	scheduleService := schedule.NewService(repositories.Schedule)
	eduService := edu.NewService(repositories.Edu)
	groupService := group.NewService(logger,
		repositories.Group,
		repositories.Member,
		repositories.User,
		scheduleService,
		repositories.Schedule,
		userService,
		eduService)

	return &Services{
		User:     userService,
		Auth:     authService,
		Schedule: scheduleService,
		Edu:      eduService,
		Group:    groupService,
	}
}
