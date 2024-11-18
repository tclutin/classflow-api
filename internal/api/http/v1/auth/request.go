package auth

type SignUpRequest struct {
	Email    string  `json:"email" binding:"required,email,max=40"`
	Password string  `json:"password" binding:"required,min=8,max=40"`
	Role     string  `json:"role" binding:"required,min=4,max=10"`
	FullName *string `json:"fullname,omitempty"`
	Telegram *string `json:"telegram,omitempty"`
}

type LogInRequest struct {
	Email    string `json:"email" binding:"required,email,max=40"`
	Password string `json:"password" binding:"required,min=8,max=40"`
}

/*
https://habr.com/ru/articles/780280/
 */