package handler

import (
	"context"
	"github.com/advanced-coder-com/go-timekeeper/internal/logs"
	"gitlab.com/tozd/go/errors"
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
		handler.processErrorResponse(ctx, http.StatusBadRequest, err, service.ErrUserInvalidInput)
		return
	}

	user, err := handler.userService.Signup(context.Background(), input)
	if err != nil {
		handler.processErrorResponse(ctx, http.StatusBadRequest, err, service.ErrUserSignUpFailed)
		return
	}

	token, err := auth.GenerateJWT(user.ID.String())
	if err != nil {
		handler.processErrorResponse(ctx, http.StatusInternalServerError, err, service.ErrUserSignUpFailed)
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
		handler.processErrorResponse(ctx, http.StatusBadRequest, err, service.ErrUserInvalidInput)
		return
	}

	user, err := handler.userService.Signin(context.Background(), input)
	if err != nil {
		handler.processErrorResponse(ctx, http.StatusUnauthorized, err, service.ErrUserSignInFailed)
		return
	}

	token, err := auth.GenerateJWT(user.ID.String())
	if err != nil {
		handler.processErrorResponse(ctx, http.StatusInternalServerError, err, service.ErrUserSignInFailed)
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
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": service.ErrUserUnauthorized.Error()})
		return
	}

	user, err := handler.userService.GetUser(ctx.Request.Context(), userID)
	if err != nil {
		handler.processErrorResponse(ctx, http.StatusNotFound, err, service.ErrGetUserFailed)
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
		handler.processErrorResponse(ctx, http.StatusBadRequest, err, service.ErrUserInvalidInput)
		return
	}

	if err := handler.userService.ChangePassword(ctx.Request.Context(), userID, input); err != nil {
		handler.processErrorResponse(ctx, http.StatusInternalServerError, err, service.ErrUserChangePasswordFailed)
		return
	}

	ctx.Status(http.StatusOK)
}

func (handler *UserHandler) DeleteCurrentUser(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": service.ErrUserUnauthorized.Error()})
		return
	}

	err := handler.userService.Delete(ctx.Request.Context(), userID)
	if err != nil {
		handler.processErrorResponse(ctx, http.StatusInternalServerError, err, service.ErrUserDeleteFailed)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

func (handler *UserHandler) processErrorResponse(ctx *gin.Context, responseCode int, err error, commonError error) {
	handler.logger.Error(err)
	var publicErr *service.PublicMessageError

	var message string
	if errors.As(err, &publicErr) {
		message = publicErr.Message
	} else {
		message = commonError.Error()
	}
	ctx.JSON(responseCode, gin.H{"error": message})
}
