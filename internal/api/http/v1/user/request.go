package user

type UpdateUserSettingsRequest struct {
	FullName             *string `json:"full_name,omitempty" binding:"required,max=40"`
	NotificationDelay    *int64  `json:"notification_delay,omitempty" binding:"required,min=5,max=60"`
	NotificationsEnabled *bool   `json:"notifications_enabled,omitempty" binding:"required"`
}
