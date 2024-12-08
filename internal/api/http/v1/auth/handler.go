package auth

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/tclutin/classflow-api/internal/api/http/middleware"
	"github.com/tclutin/classflow-api/internal/domain/auth"
	domainErr "github.com/tclutin/classflow-api/internal/domain/errors"
	"github.com/tclutin/classflow-api/internal/domain/user"
	"github.com/tclutin/classflow-api/pkg/response"
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
		authGroup.POST("/signup", middleware.JWTMiddleware(authService), middleware.RoleMiddleware(user.Admin), h.SignUp)
		authGroup.POST("/login", h.LogIn)
		authGroup.POST("/telegram/login", h.LogInWithTelegram)
		authGroup.POST("/telegram/signup", h.SignUpWithTelegram)
		authGroup.GET("/who", middleware.JWTMiddleware(authService), h.Who)
	}
}

// @Summary		SignUp with telegram chat id
// @Description	Создание студента с telegram chat id
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			input	body		SignUpWithTelegramRequest	true	"Создать студента"
// @Success		201		{object}	TokenResponse
// @Failure		400		{object}	response.APIError
// @Failure		409		{object}	response.APIError
// @Failure		500		{object}	response.APIError
// @Router			/auth/telegram/signup [post]
func (h *Handler) SignUpWithTelegram(c *gin.Context) {
	var request SignUpWithTelegramRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.NewAPIError(err.Error()))
		return
	}

	tokens, err := h.service.SignUpWithTelegram(c.Request.Context(), auth.SignUpWithTelegramDTO{
		TelegramChatID:   request.TelegramChatID,
		TelegramUsername: request.TelegramUsername,
		Fullname:         request.Fullname,
	})

	if err != nil {
		if errors.Is(err, domainErr.ErrUserAlreadyExists) {
			c.AbortWithStatusJSON(http.StatusConflict, response.NewAPIError(err.Error()))
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, response.NewAPIError("An error occurred on the server. Please try again later."))
		return
	}

	c.JSON(http.StatusCreated, TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

// @Summary		LogIn with telegram chat id
// @Description	Аутентификация студента
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			input	body		LogInWithTelegramRequest	true	"Аутентификация студента"
// @Success		200		{object}	TokenResponse
// @Failure		400		{object}	response.APIError
// @Failure		500		{object}	response.APIError
// @Router			/auth/telegram/login [post]
func (h *Handler) LogInWithTelegram(c *gin.Context) {
	var request LogInWithTelegramRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.NewAPIError(err.Error()))
		return
	}

	tokens, err := h.service.LogInWithTelegramRequest(c.Request.Context(), auth.LogInWithTelegramDTO{
		TelegramChatID: request.TelegramChatID,
	})

	if err != nil {
		if errors.Is(err, domainErr.ErrUserNotFound) {
			c.AbortWithStatusJSON(http.StatusBadRequest, response.NewAPIError(err.Error()))
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, response.NewAPIError("An error occurred on the server. Please try again later."))
		return
	}

	c.JSON(http.StatusOK, TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

// @Security		ApiKeyAuth
// @Summary		SignUp
// @Description	Создание нового админ пользователя
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			input	body		SignUpRequest	true	"Создать пользователя"
// @Success		201		{object}	TokenResponse
// @Failure		400		{object}	response.APIError
// @Failure		409		{object}	response.APIError
// @Failure		500		{object}	response.APIError
// @Router			/auth/signup [post]
func (h *Handler) SignUp(c *gin.Context) {
	var request SignUpRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.NewAPIError(err.Error()))
		return
	}

	tokens, err := h.service.SignUp(c.Request.Context(), auth.SignUpDTO{
		Email:    request.Email,
		Password: request.Password,
	})

	if err != nil {
		if errors.Is(err, domainErr.ErrUserAlreadyExists) {
			c.AbortWithStatusJSON(http.StatusConflict, response.NewAPIError(err.Error()))
			return

		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, response.NewAPIError("An error occurred on the server. Please try again later."))
		return
	}

	c.JSON(http.StatusCreated, TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

// @Summary		LogIn
// @Description	Аутентификация админ пользователя
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			input	body		LogInRequest	true	"Аутентификация пользователя"
// @Success		200		{object}	TokenResponse
// @Failure		400		{object}	response.APIError
// @Failure		404		{object}	response.APIError
// @Failure		500		{object}	response.APIError
// @Router			/auth/login [post]
func (h *Handler) LogIn(c *gin.Context) {
	var request LogInRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.NewAPIError(err.Error()))
		return
	}

	tokens, err := h.service.LogIn(c.Request.Context(), auth.LogInDTO{
		Email:    request.Email,
		Password: request.Password,
	})

	if err != nil {
		if errors.Is(err, domainErr.ErrWrongPassword) {
			c.AbortWithStatusJSON(http.StatusBadRequest, response.NewAPIError(err.Error()))
			return
		}

		if errors.Is(err, domainErr.ErrUserNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, response.NewAPIError(err.Error()))
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, response.NewAPIError("An error occurred on the server. Please try again later."))
		return
	}

	c.JSON(http.StatusOK, TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

// @Security		ApiKeyAuth
// @Summary		Who
// @Description	Получение информации о пользователе
// @Tags			auth
// @Accept			json
// @Produce		json
// @Success		200	{object}	UserDetailsResponse
// @Failure		401	{object}	response.APIError
// @Failure		404	{object}	response.APIError
// @Failure		500	{object}	response.APIError
// @Router			/auth/who [get]
func (h *Handler) Who(c *gin.Context) {
	value, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.NewAPIError("userID not found in context"))
		return
	}

	who, err := h.service.Who(c.Request.Context(), value.(uint64))
	if err != nil {
		if errors.Is(err, domainErr.ErrUserNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, response.NewAPIError(err.Error()))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.NewAPIError("An error occurred on the server. Please try again later."))
		return
	}

	c.JSON(http.StatusOK, UserDetailsResponse{
		UserID:               who.UserID,
		Email:                who.Email,
		Role:                 who.Role,
		FullName:             who.FullName,
		TelegramUsername:     who.TelegramUsername,
		TelegramChatID:       who.TelegramChatID,
		NotificationDelay:    who.NotificationDelay,
		NotificationsEnabled: who.NotificationsEnabled,
		CreatedAt:            who.CreatedAt,
	})
}
