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

func (r *ResumeRepository) GetResumeByID(userID, resumeID uuid.UUID) (*models.Resume, error) {
	var resume models.Resume
	if err := r.db.Preload("Skills").Preload("Experience").Preload("Education").Where("id = ? AND user_id = ?", resumeID, userID).First(&resume).Error; err != nil {
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
	// Ищем скилл среди всех, включая soft-deleted
	err := r.db.Unscoped().Where("name = ?", name).First(&skill).Error
	if err == nil {
		// Если найден и был soft-deleted, восстанавливаем
		if skill.DeletedAt.Valid {
			if err := r.db.Unscoped().Model(&skill).Update("deleted_at", nil).Error; err != nil {
				return nil, err
			}
			skill.DeletedAt.Valid = false
		}
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

func (r *ResumeRepository) GetSkillsByResumeID(resumeID uuid.UUID) ([]*models.Skill, error) {
	var skills []*models.Skill
	if err := r.db.Table("resume_skills").Select("skills.*").
		Joins("join skills on skills.id = resume_skills.skill_id").
		Where("resume_skills.resume_id = ?", resumeID).Scan(&skills).Error; err != nil {
		return nil, err
	}
	return skills, nil
}

func (r *ResumeRepository) DeleteSkillFromResume(resumeID, skillID uuid.UUID) error {
	return r.db.Table("resume_skills").Where("resume_id = ? AND skill_id = ?", resumeID, skillID).Delete(nil).Error
}

func (r *ResumeRepository) DeleteUnusedSkill(skillID uuid.UUID) error {
	var count int64
	if err := r.db.Table("resume_skills").Where("skill_id = ?", skillID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	return r.db.Delete(&models.Skill{}, "id = ?", skillID).Error
}

func (r *ResumeRepository) DeleteUnusedEdAndEx(resumeID uuid.UUID) error {
	if err := r.db.Delete(&models.Education{}, "resume_id = ?", resumeID).Error; err != nil {
		return err
	}
	return r.db.Delete(&models.Experience{}, "resume_id = ?", resumeID).Error
}
func (r *ResumeRepository) DeleteUnusedMatching(resumeID uuid.UUID) error {
	return r.db.Delete(&models.MatchingResult{}, "resume_id = ?", resumeID).Error
}

func (r *ResumeRepository) DeleteResumeFile(resumeID uuid.UUID) error {
	return r.db.Delete(&models.ResumeFile{}, "resume_id = ?", resumeID).Error
}

func (r *ResumeRepository) DeleteResume(resumeID uuid.UUID) error {
	return r.db.Delete(&models.Resume{}, "id = ?", resumeID).Error
}
