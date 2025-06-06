package router

import (
	"github.com/advanced-coder-com/go-timekeeper/internal/handler"
	"github.com/advanced-coder-com/go-timekeeper/internal/middleware"
	"github.com/gin-gonic/gin"
)

func setupUserRoutes(engine *gin.Engine) {
	userHandler := handler.NewUserHandler()
	user := engine.Group("/api/user")
	{
		user.POST("/signup", userHandler.Signup)
		user.POST("/signin", userHandler.Signin)
		user.GET("/profile", middleware.AuthRequired(), userHandler.Profile)
		user.DELETE("/delete", middleware.AuthRequired(), userHandler.DeleteCurrentUser)
		user.PATCH("/change-password", middleware.AuthRequired(), userHandler.ChangePassword)
	}
}
