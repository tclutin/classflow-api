package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tclutin/classflow-api/internal/domain/user"
	"log/slog"
)

type UserRepository struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func NewUserRepository(pool *pgxpool.Pool, logger *slog.Logger) *UserRepository {
	return &UserRepository{
		pool:   pool,
		logger: logger,
	}
}

func (u *UserRepository) Create(ctx context.Context, user user.User) (uint64, error) {
	sql := `INSERT INTO public.users (email, password_hash, role, fullname, telegram_username, telegram_chat, notification_delay, notifications_enabled, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING user_id`

	row := u.pool.QueryRow(
		ctx,
		sql,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.FullName,
		user.TelegramUsername,
		user.TelegramChatID,
		user.NotificationDelay,
		user.NotificationsEnabled,
		user.CreatedAt)

	var userID uint64

	if err := row.Scan(&userID); err != nil {
		u.logger.Error("Failed to create user",
			"error", err,
			"email", user.Email,
			"role", user.Role,
		)
		return 0, err
	}

	return userID, nil
}

func (u *UserRepository) Update(ctx context.Context, user user.User) error {
	sql := `
		UPDATE
			public.users
		SET
		    email = $1,
		    password_hash = $2,
		    role = $3,
		    fullname = $4,
		    telegram_chat = $5,
		    telegram_username = $6,
		    notification_delay = $7,
		    notifications_enabled = $8
		WHERE
		    user_id = $9
	`

	_, err := u.pool.Exec(
		ctx,
		sql,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.FullName,
		user.TelegramChatID,
		user.TelegramUsername,
		user.NotificationDelay,
		user.NotificationsEnabled,
		user.UserID)

	if err != nil {
		u.logger.Error("Failed to update user",
			"error", err,
			"userID", user.UserID,
		)
		return err
	}

	return nil
}

func (u *UserRepository) UpdateTx(ctx context.Context, tx pgx.Tx, user user.User) error {
	sql := `
		UPDATE
			public.users
		SET
		    email = $1,
		    password_hash = $2,
		    role = $3,
		    fullname = $4,
		    telegram_chat = $5,
		    telegram_username = $6,
		    notification_delay = $7,
		    notifications_enabled = $8
		WHERE
		    user_id = $9
	`

	_, err := tx.Exec(
		ctx,
		sql,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.FullName,
		user.TelegramChatID,
		user.TelegramUsername,
		user.NotificationDelay,
		user.NotificationsEnabled,
		user.UserID)

	if err != nil {
		u.logger.Error("Failed to update user",
			"error", err,
			"userID", user.UserID,
		)
		return err
	}

	return nil
}

func (u *UserRepository) GetById(ctx context.Context, userID uint64) (user.User, error) {
	sql := `SELECT * FROM public.users WHERE user_id = $1`

	row := u.pool.QueryRow(ctx, sql, userID)

	var usr user.User
	err := row.Scan(
		&usr.UserID,
		&usr.Email,
		&usr.PasswordHash,
		&usr.Role,
		&usr.FullName,
		&usr.TelegramUsername,
		&usr.TelegramChatID,
		&usr.NotificationDelay,
		&usr.NotificationsEnabled,
		&usr.CreatedAt,
	)

	if err != nil {
		u.logger.Error("Failed to get user by ID",
			"error", err,
			"userID", userID,
		)
		return usr, err
	}

	return usr, nil
}

func (u *UserRepository) GetByEmail(ctx context.Context, email string) (user.User, error) {
	sql := `SELECT * FROM public.users WHERE email = $1`

	row := u.pool.QueryRow(ctx, sql, email)

	var usr user.User
	err := row.Scan(
		&usr.UserID,
		&usr.Email,
		&usr.PasswordHash,
		&usr.Role,
		&usr.FullName,
		&usr.TelegramUsername,
		&usr.TelegramChatID,
		&usr.NotificationDelay,
		&usr.NotificationsEnabled,
		&usr.CreatedAt,
	)

	if err != nil {
		u.logger.Error("Failed to retrieve user by email",
			"error", err,
			"email", email,
		)
		return usr, err
	}

	return usr, nil
}

func (u *UserRepository) GetByTelegramChatId(ctx context.Context, telegramChatID int64) (user.User, error) {
	sql := `SELECT * FROM public.users WHERE telegram_chat = $1`

	row := u.pool.QueryRow(ctx, sql, telegramChatID)

	var usr user.User
	err := row.Scan(
		&usr.UserID,
		&usr.Email,
		&usr.PasswordHash,
		&usr.Role,
		&usr.FullName,
		&usr.TelegramUsername,
		&usr.TelegramChatID,
		&usr.NotificationDelay,
		&usr.NotificationsEnabled,
		&usr.CreatedAt)

	if err != nil {
		u.logger.Error("Failed to retrieve user by Telegram Chat ID",
			"error", err,
			"telegramChatID", telegramChatID,
		)
		return usr, err
	}

	return usr, nil
}
