package router

import (
	"github.com/advanced-coder-com/go-timekeeper/internal/handler"
	"github.com/advanced-coder-com/go-timekeeper/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, userHandler *handler.UserHandler) {
	// User API
	api := r.Group("/api/user")
	{
		api.POST("/signup", userHandler.Signup)
		api.POST("/signin", userHandler.Signin)
		api.GET("/profile", middleware.AuthRequired(), userHandler.Profile)
		api.DELETE("/delete", middleware.AuthRequired(), userHandler.DeleteCurrentUser)
		api.PATCH("/change-password", middleware.AuthRequired(), userHandler.ChangePassword)
	}
}
