package repository

import (
	"CVMatch/internal/models"

	"gorm.io/gorm"
)

type ResumeRepository struct {
	db *gorm.DB
}

func NewResumeRepository(db *gorm.DB) *ResumeRepository {
	return &ResumeRepository{
		db: db,
	}
}

func (r *ResumeRepository) Create(resume *models.Resume) error {
	return r.db.Create(resume).Error
}
