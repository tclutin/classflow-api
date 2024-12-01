package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tclutin/classflow-api/internal/domain/user"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (u *UserRepository) Create(ctx context.Context, user user.User) (uint64, error) {
	sql := `INSERT INTO public.users (email, password_hash, role, fullname, telegram_chat, notifications_enabled, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING user_id`

	row := u.pool.QueryRow(
		ctx,
		sql,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.FullName,
		user.TelegramChatID,
		user.NotificationsEnabled,
		user.CreatedAt)

	var userID uint64

	if err := row.Scan(&userID); err != nil {
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
		    notifications_enabled = $6
		WHERE
		    user_id = $7
	`

	_, err := u.pool.Exec(
		ctx,
		sql,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.FullName,
		user.TelegramChatID,
		user.NotificationsEnabled,
		user.UserID)

	return err
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
		&usr.TelegramChatID,
		&usr.NotificationsEnabled,
		&usr.CreatedAt,
	)

	if err != nil {
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
		&usr.TelegramChatID,
		&usr.NotificationsEnabled,
		&usr.CreatedAt,
	)

	if err != nil {
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
		&usr.TelegramChatID,
		&usr.NotificationsEnabled,
		&usr.CreatedAt)

	if err != nil {
		return usr, err
	}

	return usr, nil
}
