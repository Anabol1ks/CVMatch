// CreateResumeWithUser создаёт и сохраняет резюме, устанавливая UserID

package service

import (
	"CVMatch/internal/config"
	"CVMatch/internal/models"
	"CVMatch/internal/parser"
	"CVMatch/internal/repository"
	"CVMatch/internal/response"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ResumeService struct {
	repo *repository.ResumeRepository
	log  *zap.Logger
	cfg  *config.Config
}

func NewResumeService(repo *repository.ResumeRepository, log *zap.Logger, cfg *config.Config) *ResumeService {
	return &ResumeService{
		repo: repo,
		log:  log,
		cfg:  cfg,
	}
}

func (s *ResumeService) CreateResumeWithUser(path string, userID uuid.UUID) (*models.Resume, error) {
	llmRes, err := parser.ParseResumeWithLLM(path, s.cfg)
	if err != nil {
		s.log.Error("Failed to parse resume", zap.Error(err))
		return nil, err
	}

	llmRes = strings.TrimSpace(llmRes)
	if strings.HasPrefix(llmRes, "```") {
		llmRes = strings.TrimPrefix(llmRes, "```")
		llmRes = strings.TrimSpace(llmRes)
	}
	if strings.HasSuffix(llmRes, "```") {
		llmRes = strings.TrimSuffix(llmRes, "```")
		llmRes = strings.TrimSpace(llmRes)
	}

	var dto response.ParsedResumeDTO
	if err := json.Unmarshal([]byte(llmRes), &dto); err != nil {
		s.log.Error("Failed to unmarshal resume DTO", zap.Error(err))
		return nil, err
	}

	// Маппинг Skills
	var skills []models.Skill
	for _, skillName := range dto.Skills {
		skills = append(skills, models.Skill{Name: skillName})
	}

	// Маппинг Experience
	var experience []models.Experience
	for _, exp := range dto.Experience {
		experience = append(experience, models.Experience{
			Company:     exp.Company,
			Position:    exp.Position,
			StartDate:   exp.StartDate,
			EndDate:     exp.EndDate,
			Description: exp.Description,
		})
	}

	// Маппинг Education
	var education []models.Education
	for _, edu := range dto.Education {
		education = append(education, models.Education{
			Institution: edu.Institution,
			Degree:      edu.Degree,
			Field:       edu.Field,
			StartDate:   edu.StartDate,
			EndDate:     edu.EndDate,
		})
	}

	resume := &models.Resume{
		UserID:     userID,
		FullName:   dto.FullName,
		Email:      dto.Email,
		Phone:      dto.Phone,
		Location:   dto.Location,
		Skills:     skills,
		Experience: experience,
		Education:  education,
	}

	// Сохраняем резюме в базу данных
	if err := s.repo.Create(resume); err != nil {
		s.log.Error("Failed to save resume", zap.Error(err))
		return nil, err
	}

	return resume, nil
}
