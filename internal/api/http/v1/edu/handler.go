package edu

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/tclutin/classflow-api/internal/domain/edu"
	"net/http"
	"strconv"
)

type Service interface {
	GetAllFaculties(ctx context.Context) ([]edu.Faculty, error)
	GetAllProgramsByFacultyId(ctx context.Context, facultyID uint64) ([]edu.Program, error)
	GetAllTypesOfSubject(ctx context.Context) ([]edu.TypeOfSubject, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) Bind(router *gin.RouterGroup) {
	eduGroup := router.Group("/edu")
	{
		eduGroup.GET("/types_of_subject", h.GetAllTypesOfSubject)
		eduGroup.GET("/faculties", h.GetAllFaculties)
		eduGroup.GET("/faculties/:faculty_id/programs", h.GetProgramsByFacultyId)
	}
}

func (h *Handler) GetAllTypesOfSubject(c *gin.Context) {
	types, err := h.service.GetAllTypesOfSubject(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, EntitiesToTypesOfSubjectResponse(types))
}

func (h *Handler) GetAllFaculties(c *gin.Context) {
	faculties, err := h.service.GetAllFaculties(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, EntitiesToFacultiesResponse(faculties))
}

func (h *Handler) GetProgramsByFacultyId(c *gin.Context) {
	facultyID := c.Param("faculty_id")

	parseUint, err := strconv.ParseUint(facultyID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	programs, err := h.service.GetAllProgramsByFacultyId(c.Request.Context(), parseUint)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, EntitiesToProgramsResponse(programs))
}
