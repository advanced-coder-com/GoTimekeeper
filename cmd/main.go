package main

import (
	"fmt"
	"github.com/advanced-coder-com/go-timekeeper/internal/db"
	"github.com/advanced-coder-com/go-timekeeper/internal/logs"
	"github.com/advanced-coder-com/go-timekeeper/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const ENV = ".env"

func initConfig() {
	viper.SetConfigFile(ENV)
	viper.AutomaticEnv()
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		logger.Error("No .env file found or error reading it: %v", err)
	}
}

var logger logs.Logger

func main() {
	logger := logs.Get()
	initConfig()

	port := viper.GetString("APP_PORT")
	if port == "" {
		logger.Fatal("No APP_PORT environment variable found")
		panic("No APP_PORT environment variable found")
	}

	db.Init()

	engine := gin.Default()
	router.SetupRoutes(engine)
	logger.Info("🚀 Starting server on port %s...", port)
	if err := engine.Run(":" + port); err != nil {
		logger.Fatal("Server failed: %v", err)
		if viper.GetString("DEBUG") == "true" {
			fmt.Printf("Server failed: %v\n", err)
		}
	}
}
