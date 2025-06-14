package handler

import (
	"github.com/advanced-coder-com/go-timekeeper/internal/model"
	"github.com/advanced-coder-com/go-timekeeper/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{service: service.NewTaskService()}
}

func (taskHandler *TaskHandler) Create(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	var input service.CreateTaskInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": service.ErrTaskInvalidInput.Error()})
		return
	}

	if input.Status != "" && !model.IsValidTaskStatus(input.Status) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": service.ErrTaskInvalidInputStatus.Error()})
		return
	}

	task, err := taskHandler.service.Create(ctx.Request.Context(), userID, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": service.ErrTaskCreateFailed.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, task)
}

func (taskHandler *TaskHandler) ListAll(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	tasks, err := taskHandler.service.GetAllByUser(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": service.ErrTaskListFailed.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}

func (taskHandler *TaskHandler) ListActive(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	tasks, err := taskHandler.service.GetAllActiveByUser(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": service.ErrTaskListFailed.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}

func (taskHandler *TaskHandler) GetByID(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	taskID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": service.ErrTaskInvalidInput.Error()})
		return
	}

	task, err := taskHandler.service.GetByID(ctx.Request.Context(), taskID, userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": service.ErrTaskNotFound.Error()})
		return
	}

	ctx.JSON(http.StatusOK, task)
}

func (taskHandler *TaskHandler) Update(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	taskID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": service.ErrTaskInvalidInput.Error()})
		return
	}

	var input service.UpdateTaskInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := taskHandler.service.Update(ctx.Request.Context(), taskID, userID, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": service.ErrTaskUpdateFailed.Error()})
		return
	}

	ctx.JSON(http.StatusOK, task)
}

func (taskHandler *TaskHandler) Delete(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	taskID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": service.ErrTaskInvalidInput.Error()})
		return
	}

	if err := taskHandler.service.Delete(ctx.Request.Context(), taskID, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": service.ErrTaskDeleteFailed.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (taskHandler *TaskHandler) Start(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	taskID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": service.ErrTaskInvalidInput.Error()})
		return
	}
	if err := taskHandler.service.Start(ctx.Request.Context(), taskID, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": service.ErrTaskStartFailed.Error()})
		return
	}
	ctx.Status(http.StatusOK)
}

func (taskHandler *TaskHandler) Stop(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	taskID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": service.ErrTaskInvalidInput.Error()})
		return
	}
	if err := taskHandler.service.Stop(ctx.Request.Context(), taskID, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": service.ErrTaskStopFailed.Error()})
		return
	}
	ctx.Status(http.StatusOK)
}

func (taskHandler *TaskHandler) StopAll(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	if err := taskHandler.service.StopAll(ctx.Request.Context(), userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": service.ErrTaskStopFailed.Error()})
		return
	}
	ctx.Status(http.StatusOK)
}

func (taskHandler *TaskHandler) Close(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	taskID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": service.ErrTaskInvalidInput.Error()})
		return
	}
	if err := taskHandler.service.Stop(ctx.Request.Context(), taskID, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": service.ErrTaskUpdateFailed.Error()})
		return
	}
	ctx.Status(http.StatusOK)
}
