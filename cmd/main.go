package main

import (
	"CVMatch/internal/config"
	"CVMatch/internal/logger"
	"CVMatch/internal/router"
	"CVMatch/internal/storage"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	_ = godotenv.Load()
	isDev := os.Getenv("ENV") == "development"
	if err := logger.Init(isDev); err != nil {
		panic(err)
	}

	defer logger.Sync()

	log := logger.L()

	cfg := config.Load(log)
	db := storage.ConnectDB(&cfg.DB, log)
	storage.Migrate(db, log)

	r := router.Router(db, log, cfg)
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}
