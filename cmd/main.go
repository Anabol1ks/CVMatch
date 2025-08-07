package main

import (
	_ "CVMatch/docs"
	"CVMatch/internal/config"
	"CVMatch/internal/handlers"
	"CVMatch/internal/logger"
	"CVMatch/internal/parser"
	"CVMatch/internal/repository"
	"CVMatch/internal/router"
	"CVMatch/internal/service"
	"CVMatch/internal/storage"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// @Title CVMatch API
// @Version 1.0
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, log, cfg)
	userHandler := handlers.NewUserHandler(userService)

	resumeRepo := repository.NewResumeRepository(db)
	resumeParser := parser.YandexResumeParser{}
	resumeService := service.NewResumeService(resumeRepo, log, cfg, resumeParser)
	resumeHandler := handlers.NewResumeHandler(resumeService)

	handlers := &router.Handlers{
		User:   userHandler,
		Resume: resumeHandler,
	}

	r := router.Router(db, log, cfg, handlers)
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}
