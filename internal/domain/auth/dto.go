package auth

type TokenDTO struct {
	AccessToken  string
	RefreshToken string
}

type SignUpDTO struct {
	Email    string
	Password string
	Role     string
	FullName *string
	Telegram *string
}

type LogInDTO struct {
	Email    string
	Password string
}
