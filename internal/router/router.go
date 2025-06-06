package router

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(engine *gin.Engine) {
	// User API
	setupUserRoutes(engine)

	//Project API
	setupProjectRoutes(engine)
}
