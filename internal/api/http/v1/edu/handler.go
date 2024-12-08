package edu

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/tclutin/classflow-api/internal/api/http/middleware"
	"github.com/tclutin/classflow-api/internal/domain/auth"
	"github.com/tclutin/classflow-api/internal/domain/edu"
	"github.com/tclutin/classflow-api/pkg/response"
	"net/http"
	"strconv"
)

type Service interface {
	GetAllFaculties(ctx context.Context) ([]edu.Faculty, error)
	GetAllProgramsByFacultyId(ctx context.Context, facultyID uint64) ([]edu.Program, error)
	GetAllBuildings(ctx context.Context) ([]edu.Building, error)
	GetAllTypesOfSubject(ctx context.Context) ([]edu.TypeOfSubject, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) Bind(router *gin.RouterGroup, authService *auth.Service) {
	eduGroup := router.Group("/edu", middleware.JWTMiddleware(authService))
	{
		eduGroup.GET("/buildings", h.GetAllBuildings)
		eduGroup.GET("/types_of_subject", h.GetAllTypesOfSubject)
		eduGroup.GET("/faculties", h.GetAllFaculties)
		eduGroup.GET("/faculties/:faculty_id/programs", h.GetProgramsByFacultyId)
	}
}

// @Security		ApiKeyAuth
// @Summary		GetAllBuildings
// @Description	Получить список корпусов
// @Tags			edu
// @Accept			json
// @Produce		json
// @Success		200	{array}		BuildingResponse
// @Failure		500	{object}	response.APIError
// @Router			/edu/buildings [get]
func (h *Handler) GetAllBuildings(c *gin.Context) {
	buildings, err := h.service.GetAllBuildings(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.NewAPIError("An error occurred on the server. Please try again later."))
		return
	}

	c.JSON(http.StatusOK, EntitiesToBuildingsResponse(buildings))
}

// @Security		ApiKeyAuth
// @Summary		GetAllTypesOfSubject
// @Description	Получить список типов всех предметов
// @Tags			edu
// @Accept			json
// @Produce		json
// @Success		200	{array}		TypeOfSubjectResponse
// @Failure		500	{object}	response.APIError
// @Router			/edu/types_of_subject [get]
func (h *Handler) GetAllTypesOfSubject(c *gin.Context) {
	types, err := h.service.GetAllTypesOfSubject(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.NewAPIError("An error occurred on the server. Please try again later."))
		return
	}

	c.JSON(http.StatusOK, EntitiesToTypesOfSubjectResponse(types))
}

// @Security		ApiKeyAuth
// @Summary		GetAllFaculties
// @Description	Получить список всех факультетов
// @Tags			edu
// @Accept			json
// @Produce		json
// @Success		200	{array}		FacultyResponse
// @Failure		500	{object}	response.APIError
// @Router			/edu/faculties [get]
func (h *Handler) GetAllFaculties(c *gin.Context) {
	faculties, err := h.service.GetAllFaculties(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.NewAPIError("An error occurred on the server. Please try again later."))
		return
	}

	c.JSON(http.StatusOK, EntitiesToFacultiesResponse(faculties))
}

// @Security		ApiKeyAuth
// @Summary		GetProgramsByFacultyId
// @Description	Получить всех программ факультета
// @Tags			edu
// @Accept			json
// @Produce		json
// @Param			faculty_id	path		string	true	"faculty ID"
// @Success		200			{array}		ProgramResponse
// @Failure		400			{object}	response.APIError
// @Failure		500			{object}	response.APIError
// @Router			/edu/faculties/{faculty_id}/programs [get]
func (h *Handler) GetProgramsByFacultyId(c *gin.Context) {
	facultyID := c.Param("faculty_id")

	parseUint, err := strconv.ParseUint(facultyID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.NewAPIError(err.Error()))
		return
	}

	programs, err := h.service.GetAllProgramsByFacultyId(c.Request.Context(), parseUint)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.NewAPIError("An error occurred on the server. Please try again later."))
		return
	}

	c.JSON(http.StatusOK, EntitiesToProgramsResponse(programs))
}
