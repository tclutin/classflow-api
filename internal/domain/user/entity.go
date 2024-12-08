package user

import (
	"time"
)

const (
	Admin   = "admin"
	Leader  = "leader"
	Student = "student"
)

type User struct {
	UserID               uint64
	Email                *string
	PasswordHash         *string
	Role                 string
	FullName             *string
	TelegramUsername     *string
	TelegramChatID       *int64
	NotificationDelay    *int64
	NotificationsEnabled *bool
	CreatedAt            time.Time
}
