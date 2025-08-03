package router

import (
	"CVMatch/internal/config"
	"CVMatch/internal/handlers"
	"CVMatch/internal/middleware"

	"github.com/gin-contrib/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Handlers struct {
	User   *handlers.UserHandler
	Resume *handlers.ResumeHandler
}

func Router(db *gorm.DB, log *zap.Logger, cfg *config.Config, handlers *Handlers) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Static("/uploads", "./uploads")

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	auth := r.Group("/auth")
	{
		auth.POST("/register", handlers.User.RegisterHandler)
		auth.POST("/login", handlers.User.LoginHandler)
		auth.POST("/refresh", handlers.User.RefreshHandler)
	}

	resume := r.Group("/resumes", middleware.JWTAuth(&cfg.JWT))
	{
		resume.POST("/upload", handlers.Resume.UploadResumeHandler)
		resume.GET("/list", handlers.Resume.ListResumesHandler)
		resume.GET("/:id", handlers.Resume.GetResumeHandler)
		resume.DELETE("/:id", handlers.Resume.DeleteResumeHandler)
	}

	r.GET("/profile", middleware.JWTAuth(&cfg.JWT), handlers.User.ProfileHandler)

	// r.POST("/upload", handlers.User.UploadResumeHandler)
	return r
}
