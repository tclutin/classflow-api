package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tclutin/classflow-api/internal/api"
	"github.com/tclutin/classflow-api/internal/config"
	"github.com/tclutin/classflow-api/internal/domain"
	"github.com/tclutin/classflow-api/internal/migration"
	"github.com/tclutin/classflow-api/internal/repository"
	"github.com/tclutin/classflow-api/pkg/client/postgresql"
	"github.com/tclutin/classflow-api/pkg/jwt"
	"github.com/tclutin/classflow-api/pkg/logger"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"

	"syscall"
)

type App struct {
	server *http.Server
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func NewApp() *App {
	cfg := config.MustLoad()

	appLogger := logger.New(cfg.Environment, "logs/app.log")

	dsn := fmt.Sprintf(
		"postgresql://%v:%v@%v:%v/%v",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DbName)

	postgres := postgresql.NewPool(context.Background(), dsn)

	migrator := migration.New(postgres, appLogger)
	migrator.Init(context.Background(), cfg.Admin.Email, cfg.Admin.Password)

	repositories := repository.NewRepositories(postgres, appLogger)

	jwtManager := jwt.MustLoadTokenManager(cfg.JWT.Secret)

	services := domain.NewServices(appLogger, jwtManager, repositories, cfg)

	router := api.NewRouter(services, cfg)

	appServer := &http.Server{
		Addr:    net.JoinHostPort(cfg.HTTPServer.Address, cfg.HTTPServer.Port),
		Handler: router,
	}

	return &App{
		server: appServer,
		pool:   postgres,
		logger: appLogger,
	}
}

func (app *App) Run(ctx context.Context) {
	app.logger.Info("Starting application...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		app.logger.Info("Server is starting...")
		err := app.server.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			app.logger.Error("Server stopped with error", "error", err)
			os.Exit(1)
		}
	}()

	app.logger.Info("Server started successfully")

	<-quit
	app.logger.Info("Received shutdown signal, stopping application...")
	app.Stop(ctx)
}

func (app *App) Stop(ctx context.Context) {
	app.logger.Info("Shutting down app...")

	app.pool.Close()

	if err := app.server.Shutdown(ctx); err != nil {
		app.logger.Error("Error during server shutdown", "error", err)
		os.Exit(1)
	}

	app.logger.Info("Application shutdown completed")
}
