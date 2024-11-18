package user

import (
	"time"
)

type User struct {
	UserID       uint64    `db:"user_id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	Role         string    `db:"role"`
	FullName     *string   `db:"fullname"`
	Telegram     *string   `db:"telegram"`
	CreatedAt    time.Time `db:"created_at"`
}
