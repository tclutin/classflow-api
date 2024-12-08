package auth

type TokenDTO struct {
	AccessToken  string
	RefreshToken string
}

type SignUpDTO struct {
	Email    string
	Password string
}

type LogInDTO struct {
	Email    string
	Password string
}

type SignUpWithTelegramDTO struct {
	TelegramChatID   int64
	TelegramUsername string
	Fullname         string
}

type LogInWithTelegramDTO struct {
	TelegramChatID int64
}
