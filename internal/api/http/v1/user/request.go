package user

type UpdateUserSettingsRequest struct {
	FullName             *string `json:"full_name" binding:"omitempty,max=40"`
	NotificationDelay    *int64  `json:"notification_delay" binding:"omitempty,min=5,max=60"`
	NotificationsEnabled *bool   `json:"notifications_enabled" binding:"omitempty"`
}
