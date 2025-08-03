package main

import (
	"CVMatch/internal/config"
	"CVMatch/internal/logger"
	"CVMatch/internal/repository"
	"CVMatch/internal/service"
	"CVMatch/internal/storage"
	"fmt"
	"os"

	"github.com/google/uuid"
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

	repo := repository.NewResumeRepository(db)
	service := service.NewResumeService(repo, log, cfg)
	fmt.Println(service.CreateResumeWithUser("./uploads/resume.pdf", uuid.New()))

}
