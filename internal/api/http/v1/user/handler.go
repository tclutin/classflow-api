package user

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
	UpdatePartial(ctx context.Context, dto user.PartialUpdateUserDTO, userID uint64) error
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
	userGroup := router.Group("/users", middleware.JWTMiddleware(authService), middleware.RoleMiddleware(user.Student, user.Leader))
	{
		userGroup.PATCH("/settings", h.UpdateSettings)
	}
}

// @Security		ApiKeyAuth
// @Summary		UpdateSettings
// @Description	Обновление настроек студента
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			input	body		UpdateUserSettingsRequest	false	"Update a user's account"
// @Success		200			{string}	string
// @Failure		401	{object}	response.APIError
// @Failure		404	{object}	response.APIError
// @Failure		500	{object}	response.APIError
// @Router			/users/settings [patch]
func (h *Handler) UpdateSettings(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.NewAPIError("userID not found in context"))
		return
	}

	var request UpdateUserSettingsRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.NewAPIError(err.Error()))
		return
	}

	err := h.service.UpdatePartial(c.Request.Context(), user.PartialUpdateUserDTO{
		FullName:             request.FullName,
		NotificationDelay:    request.NotificationDelay,
		NotificationsEnabled: request.NotificationsEnabled,
	}, userID.(uint64))

	if err != nil {
		if errors.Is(err, domainErr.ErrUserNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, response.NewAPIError(err.Error()))
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, response.NewAPIError("An error occurred on the server. Please try again later."))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
