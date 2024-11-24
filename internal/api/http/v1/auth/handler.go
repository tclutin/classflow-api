package auth

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	_ "github.com/tclutin/classflow-api/docs"
	"github.com/tclutin/classflow-api/internal/api/http/middleware"
	"github.com/tclutin/classflow-api/internal/domain/auth"
	domainErr "github.com/tclutin/classflow-api/internal/domain/errors"
	"github.com/tclutin/classflow-api/internal/domain/user"
	"net/http"
)

type Service interface {
	SignUp(ctx context.Context, dto auth.SignUpDTO) (auth.TokenDTO, error)
	LogIn(ctx context.Context, dto auth.LogInDTO) (auth.TokenDTO, error)
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
		authGroup.GET("/who", middleware.JWTMiddleware(authService), h.Who)
	}
}

func (h *Handler) SignUp(c *gin.Context) {
	var request SignUpRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := h.service.SignUp(c.Request.Context(), auth.SignUpDTO{
		Email:    request.Email,
		Password: request.Password,
		Role:     request.Role,
		FullName: request.FullName,
		Telegram: request.Telegram,
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
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, UserDetailsResponse{
		UserID:       who.UserID,
		Email:        who.Email,
		PasswordHash: who.PasswordHash,
		Role:         who.Role,
		FullName:     who.FullName,
		Telegram:     who.Telegram,
		CreatedAt:    who.CreatedAt,
	})
}
