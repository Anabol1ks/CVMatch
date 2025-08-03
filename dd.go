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

	// repo := repository.NewResumeRepository(db)

	// fmt.Println(repo.GetResFileURL(uuid.MustParse("")))
	// service := service.NewResumeService(repo, log, cfg)
	// fmt.Println(service.GetResumeFileURL(uuid.MustParse("")))
	// fmt.Println(repo.GetResumeByID(uuid.MustParse("1f05c867-7895-4d5c-b216-e9a103435a8e")))
	// fmt.Println(service.CreateResumeWithUser("./uploads/resume.pdf", uuid.New()))
}
