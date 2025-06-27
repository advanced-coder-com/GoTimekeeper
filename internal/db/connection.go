package db

import (
	"fmt"
	"github.com/advanced-coder-com/go-timekeeper/internal/logs"
	"log"
	"sync"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	instance *gorm.DB
	once     sync.Once
	logger   logs.Logger
)

func Init() {
	logger = logs.Get()
	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			viper.GetString("DB_HOST"),
			viper.GetString("DB_USER"),
			viper.GetString("DB_PASSWORD"),
			viper.GetString("DB_NAME"),
			viper.GetString("DB_PORT"),
		)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			logger.Fatal("❌ Failed to connect to database: %v", err)
		}

		instance = db
		logger.Info("✅ Database connection established")
	})
}

func Get() *gorm.DB {
	if instance == nil {
		log.Fatal("❌ DB not initialized — call db.Init() first")
	}
	return instance
}
