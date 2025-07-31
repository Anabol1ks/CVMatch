package storage

import (
	"CVMatch/internal/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB, log *zap.Logger) {
	if err := db.AutoMigrate(
		&models.User{},
		&models.Resume{},
		&models.ResumeFile{},
		&models.Skill{},
		&models.Experience{},
		&models.Education{},
		&models.Vacancy{},
		&models.MatchingResult{},
	); err != nil {
		log.Fatal("Ошибка миграции базы данных", zap.Error(err))
	}
	log.Info("Миграция базы данных прошла успешно")
}
