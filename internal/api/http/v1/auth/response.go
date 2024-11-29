package auth

import (
	"time"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserDetailsResponse struct {
	UserID               uint64    `json:"user_id"`
	Email                *string   `json:"email"`
	PasswordHash         *string   `json:"password_hash"`
	Role                 string    `json:"role"`
	FullName             *string   `json:"full_name"`
	TelegramChatID       *int64    `json:"telegram"`
	NotificationsEnabled *bool     `json:"notifications_enabled"`
	CreatedAt            time.Time `json:"created_at"`
}
