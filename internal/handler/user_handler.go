package handler

import (
	"context"
	"github.com/advanced-coder-com/go-timekeeper/internal/logs"
	"net/http"

	"github.com/advanced-coder-com/go-timekeeper/internal/auth"
	"github.com/advanced-coder-com/go-timekeeper/internal/service"
	"github.com/gin-gonic/gin"
)

const userHandlerErrorPrefix = "UserHandler"

type UserHandler struct {
	userService *service.UserService
	logger      logs.Logger
}

func NewUserHandler() *UserHandler {
	userService := service.NewUserService()

	return &UserHandler{
		userService: userService,
		logger:      logs.Get(),
	}
}

func (handler *UserHandler) Signup(ctx *gin.Context) {
	var input service.UserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		handler.logger.LogError(userHandlerErrorPrefix, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": service.ErrUserInvalidInput.Error()})
		return
	}

	user, err := handler.userService.Signup(context.Background(), input)
	if err != nil {
		handler.logger.LogError(userHandlerErrorPrefix, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": service.ErrUserSignUpFailed.Error()})
		return
	}

	token, err := auth.GenerateJWT(user.ID.String())
	if err != nil {
		handler.logger.LogError(userHandlerErrorPrefix, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": service.ErrUserSignUpFailed.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":    user.ID,
		"email": user.Email,
		"token": token,
	})
}

func (handler *UserHandler) Signin(ctx *gin.Context) {
	var input service.UserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		handler.logger.LogError(userHandlerErrorPrefix, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": service.ErrUserInvalidInput.Error()})
		return
	}

	user, err := handler.userService.Signin(context.Background(), input)
	if err != nil {
		handler.logger.LogError(userHandlerErrorPrefix, err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": service.ErrUserSignInFailed.Error()})
		return
	}

	token, err := auth.GenerateJWT(user.ID.String())
	if err != nil {
		handler.logger.LogError(userHandlerErrorPrefix, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": service.ErrUserSignInFailed.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"email": user.Email,
		"token": token,
	})
}

func (handler *UserHandler) Profile(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		handler.logger.LogError(userHandlerErrorPrefix, service.ErrUserUnauthorized)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": service.ErrUserUnauthorized.Error()})
		return
	}

	user, err := handler.userService.GetUser(ctx.Request.Context(), userID)
	if err != nil {
		handler.logger.LogError(userHandlerErrorPrefix, err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": service.ErrGetUserFailed.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"email": user.Email,
	})
}

func (handler *UserHandler) ChangePassword(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	var input service.ChangePasswordInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		handler.logger.LogError(userHandlerErrorPrefix, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": service.ErrUserInvalidInput.Error()})
		return
	}

	if err := handler.userService.ChangePassword(ctx.Request.Context(), userID, input); err != nil {
		handler.logger.LogError(userHandlerErrorPrefix, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": service.ErrUserChangePasswordFailed.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (handler *UserHandler) DeleteCurrentUser(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		handler.logger.LogError(userHandlerErrorPrefix, service.ErrUserUnauthorized)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": service.ErrUserUnauthorized.Error()})
		return
	}

	err := handler.userService.Delete(ctx.Request.Context(), userID)
	if err != nil {
		handler.logger.LogError(userHandlerErrorPrefix, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": service.ErrUserDeleteFailed.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
