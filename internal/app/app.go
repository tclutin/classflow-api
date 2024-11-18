package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/tclutin/classflow-api/internal/api"
	"github.com/tclutin/classflow-api/internal/config"
	"github.com/tclutin/classflow-api/internal/domain"
	"github.com/tclutin/classflow-api/internal/repository"
	"github.com/tclutin/classflow-api/pkg/client/postgresql"
	"github.com/tclutin/classflow-api/pkg/jwt"
	"net"
	"net/http"
	"os"
	"os/signal"

	"syscall"
)

type App struct {
	server *http.Server
}

func NewApp() *App {
	cfg := config.MustLoad()

	dsn := fmt.Sprintf(
		"postgresql://%v:%v@%v:%v/%v",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DbName)

	postgres := postgresql.NewPool(context.Background(), dsn)

	repositories := repository.NewRepositories(postgres)

	jwtManager := jwt.MustLoadTokenManager(cfg.JWT.Secret)

	services := domain.NewServices(jwtManager, repositories, cfg)

	router := api.NewRouter(services)

	appServer := &http.Server{
		Addr:    net.JoinHostPort(cfg.HTTPServer.Address, cfg.HTTPServer.Port),
		Handler: router,
	}

	return &App{
		server: appServer,
	}
}

func (app *App) Run(ctx context.Context) {

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		err := app.server.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			os.Exit(1)
		}
	}()

	<-quit
	app.Stop(ctx)
}

func (app *App) Stop(ctx context.Context) {
	if err := app.server.Shutdown(ctx); err != nil {
		os.Exit(1)
	}
}
