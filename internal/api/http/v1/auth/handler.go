package auth

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/tclutin/classflow-api/internal/api/http/middleware"
	"github.com/tclutin/classflow-api/internal/domain/auth"
	domainErr "github.com/tclutin/classflow-api/internal/domain/errors"
	"github.com/tclutin/classflow-api/internal/domain/user"
	"net/http"
)

type Service interface {
	SignUp(ctx context.Context, dto auth.SignUpDTO) (auth.TokenDTO, error)
	LogIn(ctx context.Context, dto auth.LogInDTO) (auth.TokenDTO, error)
	SignUpWithTelegram(ctx context.Context, dto auth.SignUpWithTelegramDTO) (auth.TokenDTO, error)
	LogInWithTelegramRequest(ctx context.Context, dto auth.LogInWithTelegramDTO) (auth.TokenDTO, error)
	Who(ctx context.Context, userID uint64) (user.User, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Bind(router *gin.RouterGroup, authService *auth.Service) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/signup", h.SignUp)
		authGroup.POST("/login", h.LogIn)
		authGroup.POST("/telegram/login", h.LogInWithTelegram)
		authGroup.POST("/telegram/signup", h.SignUpWithTelegram)
		authGroup.GET("/who", middleware.JWTMiddleware(authService), h.Who)
	}
}

func (h *Handler) SignUpWithTelegram(c *gin.Context) {
	var request SignUpWithTelegramRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := h.service.SignUpWithTelegram(c.Request.Context(), auth.SignUpWithTelegramDTO{
		TelegramChatID: request.TelegramChatID,
		Fullname:       request.Fullname,
	})

	if err != nil {
		if errors.Is(err, domainErr.ErrUserAlreadyExists) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

func (h *Handler) LogInWithTelegram(c *gin.Context) {
	var request LogInWithTelegramRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := h.service.LogInWithTelegramRequest(c.Request.Context(), auth.LogInWithTelegramDTO{
		TelegramChatID: request.TelegramChatID,
	})

	if err != nil {
		if errors.Is(err, domainErr.ErrUserNotFound) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

type HTTPError struct {
	Code    int    `json:"code"`    // HTTP статус код
	Message string `json:"message"` // Сообщение об ошибке
}

// SignUp godoc
// @Summary      Register a new account
// @Description  Create a new user account with email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body     SignUpRequest  true  "Register new account"
// @Success      201  {object}  TokenResponse
// @Failure      400  {object}  HTTPError
// @Failure      409  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /auth/signup [post]
func (h *Handler) SignUp(c *gin.Context) {
	var request SignUpRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, HTTPError{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	tokens, err := h.service.SignUp(c.Request.Context(), auth.SignUpDTO{
		Email:    request.Email,
		Password: request.Password,
	})

	if err != nil {
		if errors.Is(err, domainErr.ErrUserAlreadyExists) {
			c.AbortWithStatusJSON(http.StatusConflict, HTTPError{Code: http.StatusConflict, Message: err.Error()})
			return

		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, HTTPError{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

func (h *Handler) LogIn(c *gin.Context) {
	var request LogInRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := h.service.LogIn(c.Request.Context(), auth.LogInDTO{
		Email:    request.Email,
		Password: request.Password,
	})

	if err != nil {
		if errors.Is(err, domainErr.ErrWrongPassword) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, domainErr.ErrUserNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

func (h *Handler) Who(c *gin.Context) {
	value, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	who, err := h.service.Who(c.Request.Context(), value.(uint64))
	if err != nil {
		if errors.Is(err, domainErr.ErrUserNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, UserDetailsResponse{
		UserID:               who.UserID,
		Email:                who.Email,
		PasswordHash:         who.PasswordHash,
		Role:                 who.Role,
		FullName:             who.FullName,
		TelegramChatID:       who.TelegramChatID,
		NotificationsEnabled: who.NotificationsEnabled,
		CreatedAt:            who.CreatedAt,
	})
}
