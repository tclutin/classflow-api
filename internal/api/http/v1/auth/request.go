package auth

type SignUpRequest struct {
	Email    string `json:"email" binding:"required,email,max=40"`
	Password string `json:"password" binding:"required,min=8,max=40"`
}

type LogInRequest struct {
	Email    string `json:"email" binding:"required,email,max=40"`
	Password string `json:"password" binding:"required,min=8,max=40"`
}

type SignUpWithTelegramRequest struct {
	TelegramChatID   int64  `json:"telegram_chat_id" binding:"required,gte=1"`
	TelegramUsername string `json:"telegram_username" binding:"required,max=40"`
	Fullname         string `json:"full_name" binding:"required,max=40"`
}

type LogInWithTelegramRequest struct {
	TelegramChatID int64 `json:"telegram_chat_id" binding:"required,gte=1"`
}
