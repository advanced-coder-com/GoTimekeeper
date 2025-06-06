package router

import (
	"github.com/advanced-coder-com/go-timekeeper/internal/handler"
	"github.com/advanced-coder-com/go-timekeeper/internal/middleware"
	"github.com/gin-gonic/gin"
)

func setupProjectRoutes(engine *gin.Engine) {
	projectHandler := handler.NewProjectHandler()
	projects := engine.Group("/api/projects", middleware.AuthRequired())
	{
		projects.POST("/create", projectHandler.Create)
		projects.GET("/list", projectHandler.List)
		projects.GET("/detail/:id", projectHandler.GetByID)
		projects.PATCH("/update/:id", projectHandler.Rename)
		projects.DELETE("/delete/:id", projectHandler.Delete)
	}
}
