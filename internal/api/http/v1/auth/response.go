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
	Role                 string    `json:"role"`
	FullName             *string   `json:"full_name"`
	TelegramUsername     *string   `json:"telegram_username"`
	TelegramChatID       *int64    `json:"telegram"`
	NotificationDelay    *int64    `json:"notification_delay"`
	NotificationsEnabled *bool     `json:"notifications_enabled"`
	CreatedAt            time.Time `json:"created_at"`
}
