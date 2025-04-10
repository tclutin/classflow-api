package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/tclutin/classflow-api/internal/domain/auth"
	"github.com/tclutin/classflow-api/internal/domain/group"
	"github.com/tclutin/classflow-api/internal/metric"
	"github.com/tclutin/classflow-api/pkg/response"
	"net/http"
	"strconv"
	"strings"
)

func JWTMiddleware(authService *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.NewAPIError("authorization header is required"))
			return
		}

		parts := strings.Split(token, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.NewAPIError("authorization header must be in the format 'Bearer <token>'"))
			return
		}

		user, err := authService.VerifyAndGetCredentials(c.Request.Context(), parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.NewAPIError(err.Error()))
			return
		}

		c.Set("userID", user.UserID)
		c.Set("role", user.Role)
		c.Next()
	}
}

func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := c.Get("role")
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, response.NewAPIError("role not found in context"))
			return
		}

		extractRole, ok := role.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, response.NewAPIError("role type is invalid in context"))
			return
		}

		for _, r := range roles {
			if extractRole == r {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, response.NewAPIError("you do not have permission to access this resource"))
	}
}

func CounterRequestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		metric.IncRequestCounter(c.FullPath())
	}
}

func ScheduleRequestCounterMiddleware(groupService *group.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		groupId := c.Param("group_id")

		id, err := strconv.ParseUint(groupId, 10, 64)
		if err != nil {
			c.Next()
			return
		}

		group, err := groupService.GetById(c.Request.Context(), id)
		if err != nil {
			c.Next()
			return
		}

		c.Next()

		metric.IncScheduleRequestCounter(c.FullPath(), group.ShortName)
	}
}
