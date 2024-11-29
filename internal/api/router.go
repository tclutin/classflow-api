package api

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	v1 "github.com/tclutin/classflow-api/internal/api/http/v1"
	"github.com/tclutin/classflow-api/internal/domain"
	"net/http"
)

func NewRouter(services *domain.Services) *gin.Engine {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Use(gin.Logger(), gin.Recovery())
	gin.SetMode(gin.DebugMode)

	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	root := router.Group("/api")
	{
		v1.NewHandler(services).InitAPI(root)
	}

	return router
}
