package router

import (
	"github.com/advanced-coder-com/go-timekeeper/internal/handler"
	"github.com/advanced-coder-com/go-timekeeper/internal/middleware"
	"github.com/gin-gonic/gin"
)

func setupTaskRoutes(engine *gin.Engine) {
	taskHandler := handler.NewTaskHandler()
	tasks := engine.Group("/api/tasks", middleware.AuthRequired())
	{
		tasks.POST("/create", taskHandler.Create)
		tasks.GET("/list-all", taskHandler.ListAll)
		tasks.GET("/list-active", taskHandler.ListActive)
		tasks.GET("/detail/:id", taskHandler.GetByID)
		tasks.PATCH("/update/:id", taskHandler.Update)
		tasks.DELETE("/delete/:id", taskHandler.Delete)
		tasks.GET("/start/:id", taskHandler.Start)
		tasks.GET("/stop/:id", taskHandler.Stop)
		tasks.GET("/stop-all", taskHandler.StopAll)
		tasks.GET("/close/:id", taskHandler.Close)
	}
}
