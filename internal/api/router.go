package api

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"github.com/tclutin/classflow-api/internal/config"

	v1 "github.com/tclutin/classflow-api/internal/api/http/v1"
	"github.com/tclutin/classflow-api/internal/domain"
	"net/http"
)

func NewRouter(services *domain.Services, cfg *config.Config) *gin.Engine {
	if cfg.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}

	if cfg.IsLocal() {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()

	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	root := router.Group("/api")
	{
		v1.NewHandler(services).InitAPI(root)
	}

	return router
}
