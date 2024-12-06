package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/tclutin/classflow-api/internal/api/http/v1/auth"
	"github.com/tclutin/classflow-api/internal/api/http/v1/edu"
	"github.com/tclutin/classflow-api/internal/api/http/v1/group"
	"github.com/tclutin/classflow-api/internal/api/http/v1/user"
	"github.com/tclutin/classflow-api/internal/domain"
)

type Handler struct {
	services *domain.Services
}

func NewHandler(services *domain.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitAPI(router *gin.RouterGroup) {
	apiGroup := router.Group("/v1")
	{
		user.NewHandler(h.services.User).Bind(apiGroup, h.services.Auth)
		auth.NewHandler(h.services.Auth).Bind(apiGroup, h.services.Auth)
		group.NewHandler(h.services.Group).Bind(apiGroup, h.services.Auth)
		edu.NewHandler(h.services.Edu).Bind(apiGroup, h.services.Auth)
	}
}
