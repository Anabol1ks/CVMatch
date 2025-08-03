package repository

import (
	"CVMatch/internal/models"
	"fmt"

	"github.com/google/uuid"
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

func (r *ResumeRepository) CreateFile(file *models.ResumeFile) error {
	return r.db.Create(file).Error
}

func (r *ResumeRepository) GetResumeByID(userID, id uuid.UUID) (*models.Resume, error) {
	var resume models.Resume
	if err := r.db.Preload("Skills").Preload("Experience").Preload("Education").First(&resume, id).Error; err != nil {
		return nil, err
	}
	return &resume, nil
}

func (r *ResumeRepository) GetListRes(userID uuid.UUID) (*[]models.Resume, error) {
	var resumes []models.Resume
	if err := r.db.Where("user_id = ?", userID).Find(&resumes).Error; err != nil {
		return nil, err
	}
	return &resumes, nil
}

func (r *ResumeRepository) GetResFileURL(id uuid.UUID) (string, error) {
	var fileURL string
	if err := r.db.Model(&models.ResumeFile{}).Where("resume_id = ?", id).Select("path").Scan(&fileURL).Error; err != nil {
		return "", err
	}
	return fileURL, nil
}

func (r *ResumeRepository) FirstOrCreateSkill(name string) (*models.Skill, error) {
	var skill models.Skill
	err := r.db.Where("name = ?", name).First(&skill).Error
	if err == nil {
		fmt.Println("Skill found:", skill.Name)
		return &skill, nil
	}
	if err == gorm.ErrRecordNotFound {
		skill.Name = name
		if err := r.db.Create(&skill).Error; err != nil {
			fmt.Println("Failed to create skill:", skill.Name)
			return nil, err
		}
		fmt.Println("Created skill:", skill.Name)
		return &skill, nil
	}
	return nil, err
}
