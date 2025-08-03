package repository

import (
	"CVMatch/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ResumeRepository struct {
	db *gorm.DB
}

// Возвращает *gorm.DB для прямого доступа (например, для select по именам)
func (r *ResumeRepository) DB() *gorm.DB {
	return r.db
}

// Ассоциация резюме и скиллов через many2many
func (r *ResumeRepository) AssociateSkills(resume *models.Resume, skills []*models.Skill) error {
	var rows []map[string]interface{}
	for _, skill := range skills {
		rows = append(rows, map[string]interface{}{
			"resume_id": resume.ID,
			"skill_id":  skill.ID,
		})
	}
	return r.db.Table("resume_skills").Clauses(clause.OnConflict{DoNothing: true}).Create(rows).Error
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
		return &skill, nil
	}
	if err == gorm.ErrRecordNotFound {
		skill.Name = name
		if err := r.db.Create(&skill).Error; err != nil {
			return nil, err
		}
		// После создания обязательно получить объект из базы по имени (гарантия ID)
		if err := r.db.Where("name = ?", name).First(&skill).Error; err != nil {
			return nil, err
		}
		return &skill, nil
	}
	return nil, err
}

func (r *ResumeRepository) WithTx(tx *gorm.DB) *ResumeRepository {
	return &ResumeRepository{db: tx}
}
