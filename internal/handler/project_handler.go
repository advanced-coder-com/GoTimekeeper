package handler

import (
	"net/http"
	"strings"

	"github.com/advanced-coder-com/go-timekeeper/internal/service"
	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	projectService *service.ProjectService
}

func NewProjectHandler() *ProjectHandler {
	return &ProjectHandler{
		projectService: service.NewProjectService(),
	}
}

func (projectHandler *ProjectHandler) Create(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	var input service.ProjectInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	project, err := projectHandler.projectService.Create(ctx.Request.Context(), userID, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": service.ErrProjectCreateFailed.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, project)
}

func (projectHandler *ProjectHandler) GetByID(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	projectID := ctx.Param("id")

	project, err := projectHandler.projectService.GetByID(ctx.Request.Context(), projectID, userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": service.ErrProjectNotFound.Error()})
		return
	}
	ctx.JSON(http.StatusOK, project)
}

func (projectHandler *ProjectHandler) List(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	projects, err := projectHandler.projectService.GetAllByUser(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": service.ErrProjectListFailed.Error()})
		return
	}
	ctx.JSON(http.StatusOK, projects)
}

func (projectHandler *ProjectHandler) Rename(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	projectID := ctx.Param("id")

	var payload struct {
		Name string `json:"name"`
	}
	if err := ctx.ShouldBindJSON(&payload); err != nil || strings.TrimSpace(payload.Name) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid name"})
		return
	}

	err := projectHandler.projectService.Rename(ctx.Request.Context(), projectID, userID, payload.Name)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": service.ErrProjectUpdateFailed.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "project renamed"})
}

func (projectHandler *ProjectHandler) Delete(ctx *gin.Context) {
	projectID := ctx.Param("id")
	err := projectHandler.projectService.Delete(ctx.Request.Context(), projectID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": service.ErrProjectDeleteFailed.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "project deleted"})
}
