package migration

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/tclutin/classflow-api/internal/domain/user"
	"github.com/tclutin/classflow-api/pkg/hash"
	"log/slog"
	"os"
	"time"
)

type Migration struct {
	logger *slog.Logger
	pool   *pgxpool.Pool
}

func New(pool *pgxpool.Pool, logger *slog.Logger) *Migration {
	return &Migration{
		logger: logger,
		pool:   pool,
	}
}

func (m *Migration) Init(ctx context.Context, email string, password string) {
	m.Up()
	m.CreateAdminUser(ctx, email, password)
}

// TODO: проверять на ластовую миграцию/заспидранил
func (m *Migration) Up() {
	db := stdlib.OpenDBFromPool(m.pool)

	if err := goose.SetDialect("postgres"); err != nil {
		m.logger.Error("error to set db dialect",
			"error", err,
		)
		os.Exit(1)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		m.logger.Error("error applying migrations",
			"error", err,
		)
		os.Exit(1)
	}
}

func (m *Migration) CreateAdminUser(ctx context.Context, email string, password string) {
	conn, err := m.pool.Acquire(ctx)
	if err != nil {
		m.logger.Error("error acquiring connection from pool",
			"error", err,
		)
		os.Exit(1)
	}
	defer conn.Release()

	password_hash, err := hash.NewBcryptHash(password)
	if err != nil {
		m.logger.Error("error hashing password",
			"error", err,
		)
		os.Exit(1)
	}

	var existingEmail string

	sql := `SELECT email FROM public.users WHERE email = $1`

	row := conn.QueryRow(ctx, sql, email)

	if err = row.Scan(&existingEmail); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			m.logger.Error("error querying for existing email", "error", err)
			os.Exit(1)
		}
	}
	if existingEmail != "" {
		m.logger.Warn("admin account with this email already exists", "email", existingEmail)
		return
	}

	sql = `INSERT INTO public.users (email, password_hash, role, created_at) VALUES ($1, $2, $3, $4)`
	_, err = conn.Exec(ctx, sql, email, password_hash, user.Admin, time.Now())
	if err != nil {
		m.logger.Error("error inserting user",
			"error", err,
		)
		os.Exit(1)
	}

	m.logger.Info("admin user created successfully", "email", email)

}
