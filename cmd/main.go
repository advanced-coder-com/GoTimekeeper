package main

import (
	"github.com/advanced-coder-com/go-timekeeper/internal/router"
	"log"

	"github.com/advanced-coder-com/go-timekeeper/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const ENV = ".env"

func initConfig() {
	viper.SetConfigFile(ENV)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No .env file found or error reading it: %v", err)
	}
}

func main() {
	initConfig()

	port := viper.GetString("APP_PORT")
	if port == "" {
		panic("No APP_PORT environment variable found")
	}

	db.Init()

	engine := gin.Default()
	router.SetupRoutes(engine)

	log.Printf("🚀 Starting server on port %s...", port)
	if err := engine.Run(":" + port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
