package user

type PartialUpdateUserDTO struct {
	FullName             *string
	NotificationDelay    *int64
	NotificationsEnabled *bool
}
