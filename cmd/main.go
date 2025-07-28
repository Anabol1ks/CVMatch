package main

import (
	"CVMatch/internal/config"
	"CVMatch/internal/logger"
	"CVMatch/internal/storage"
	"os"

	"github.com/joho/godotenv"
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

}
