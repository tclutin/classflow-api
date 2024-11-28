package group

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/tclutin/classflow-api/internal/api/http/middleware"
	"github.com/tclutin/classflow-api/internal/domain/auth"
	domainErr "github.com/tclutin/classflow-api/internal/domain/errors"
	"github.com/tclutin/classflow-api/internal/domain/group"
	"github.com/tclutin/classflow-api/internal/domain/schedule"
	"net/http"
	"strconv"
)

type Service interface {
	Create(ctx context.Context, dto group.CreateGroupDTO) (uint64, error)
	GetStudentGroupByUserId(ctx context.Context, userID uint64) (group.SummaryGroupDTO, error)
	GetLeaderGroupsByUserId(ctx context.Context, userID uint64) ([]group.DetailsGroupDTO, error)
	GetAllGroupsSummary(ctx context.Context, filter group.FilterDTO) ([]group.SummaryGroupDTO, error)
	JoinToGroup(ctx context.Context, code string, userID, groupID uint64) error
	LeaveFromGroup(ctx context.Context, userID uint64) error
	UploadSchedule(ctx context.Context, schedule []schedule.Schedule, groupID, userID uint64) error
	GetAllSchedulesByGroupIdAndUserId(ctx context.Context, groupID uint64) ([]schedule.DetailsScheduleDTO, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Bind(router *gin.RouterGroup, authService *auth.Service) {
	groupsGroup := router.Group("/groups")
	{
		groupsGroup.POST("", middleware.JWTMiddleware(authService), middleware.RoleMiddleware("leader"), h.Create)
		groupsGroup.GET("", h.GetAllGroupsSummary)
		groupsGroup.GET("/my", middleware.JWTMiddleware(authService), h.GetGroupForCurrentUser)
		groupsGroup.POST("/:group_id/join", middleware.JWTMiddleware(authService), middleware.RoleMiddleware("student"), h.JoinToGroup)
		groupsGroup.POST("/:group_id/schedule", middleware.JWTMiddleware(authService), middleware.RoleMiddleware("leader"), h.UploadSchedule)
		groupsGroup.POST("/leave", middleware.JWTMiddleware(authService), middleware.RoleMiddleware("student"), h.LeaveFromGroup)
		groupsGroup.GET("/:group_id/schedule", h.GetScheduleByGroupId)
	}
}

func (h *Handler) Create(c *gin.Context) {
	var request CreateGroupRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	value, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "wtf"})
		return
	}

	groupID, err := h.service.Create(c.Request.Context(), group.CreateGroupDTO{
		LeaderID:  value.(uint64),
		FacultyID: request.FacultyID,
		ProgramID: request.ProgramID,
		ShortName: request.ShortName,
	})

	if err != nil {
		if errors.Is(err, domainErr.ErrGroupAlreadyExists) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, domainErr.ErrFacultyProgramIdMismatch) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, domainErr.ErrProgramNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, domainErr.ErrFacultyNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"group_id": groupID,
	})
}

func (h *Handler) LeaveFromGroup(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "wtf"})
		return
	}

	// TOOO: error handler
	if err := h.service.LeaveFromGroup(c.Request.Context(), userID.(uint64)); err != nil {
		if errors.Is(err, domainErr.ErrGroupNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, domainErr.ErrMemberNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *Handler) JoinToGroup(c *gin.Context) {
	var request JoinToGroupRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	groupID, err := strconv.ParseUint(c.Param("group_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "wtf"})
		return
	}

	err = h.service.JoinToGroup(c.Request.Context(), request.Code, userID.(uint64), groupID)
	if err != nil {
		if errors.Is(err, domainErr.ErrGroupNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, domainErr.ErrWrongGroupCode) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, domainErr.ErrAlreadyInGroup) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *Handler) UploadSchedule(c *gin.Context) {
	var request UploadScheduleRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := request.Validate(); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	groupID, err := strconv.ParseUint(c.Param("group_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "wtf"})
		return
	}

	//TODO: need to fix double groupID
	if err = h.service.UploadSchedule(c.Request.Context(), request.TransformToEntities(groupID), groupID, userID.(uint64)); err != nil {
		if errors.Is(err, domainErr.ErrGroupAlreadyHasSchedule) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, domainErr.ErrFacultyNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, domainErr.ErrGroupNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, domainErr.ErrThisGroupDoesNotBelongToYou) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, domainErr.ErrProgramNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *Handler) GetScheduleByGroupId(c *gin.Context) {
	groupID, err := strconv.ParseUint(c.Param("group_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schedules, err := h.service.GetAllSchedulesByGroupIdAndUserId(c.Request.Context(), groupID)
	if err != nil {
		if errors.Is(err, domainErr.ErrYouArentMember) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, EntitiesToSchedulesResponse(schedules))
}

func (h *Handler) GetAllGroupsSummary(c *gin.Context) {

	program := c.DefaultQuery("program", "")
	faculty := c.DefaultQuery("faculty", "")

	groups, err := h.service.GetAllGroupsSummary(c.Request.Context(), group.FilterDTO{
		Faculty: faculty,
		Program: program,
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, EntitiesToSummaryGroupsResponse(groups))
}

func (h *Handler) GetGroupForCurrentUser(c *gin.Context) {
	value, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	role, ok := c.Get("role")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if role == "student" {
		studentGroup, err := h.service.GetStudentGroupByUserId(c.Request.Context(), value.(uint64))
		if err != nil {
			if errors.Is(err, domainErr.ErrGroupNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, nil)
				return
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, SummaryGroupResponse{
			GroupID:        studentGroup.GroupID,
			Faculty:        studentGroup.Faculty,
			Program:        studentGroup.Program,
			ShortName:      studentGroup.ShortName,
			NumberOfPeople: studentGroup.NumberOfPeople,
			ExistsSchedule: studentGroup.ExistsSchedule,
		})
		return
	}

	if role == "leader" {
		leaderGroups, err := h.service.GetLeaderGroupsByUserId(c.Request.Context(), value.(uint64))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, EntitiesToDetailsGroupsResponse(leaderGroups))
		return
	}

	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
}
