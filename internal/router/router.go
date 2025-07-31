package router

import (
	"CVMatch/internal/config"
	"CVMatch/internal/handlers"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Handlers struct {
	User *handlers.UserHandler
}

func Router(db *gorm.DB, log *zap.Logger, cfg *config.Config, handlers *Handlers) *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	auth := r.Group("/auth")
	{
		auth.POST("/register", handlers.User.RegisterHandler)
	}

	// r.POST("/upload", handlers.User.UploadResumeHandler)
	return r
}
