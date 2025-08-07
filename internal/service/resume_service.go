package service

import (
	"CVMatch/internal/config"
	"CVMatch/internal/models"
	"CVMatch/internal/parser"
	"CVMatch/internal/repository"
	"CVMatch/internal/response"
	"encoding/json"
	"os"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
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

func (s *ResumeService) CreateResumeWithUser(path string, userID uuid.UUID) (*response.ParsedResumeDTO, error) {
	llmRes, err := parser.ParseResumeWithYandex(path, s.cfg)
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

	// Открываем транзакцию
	txErr := s.repo.DB().Transaction(func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		var skills []*models.Skill
		for _, skillName := range dto.Skills {
			skill, err := txRepo.FirstOrCreateSkill(skillName)
			if err != nil {
				s.log.Error("Failed to find or create skill", zap.String("skill", skillName), zap.Error(err))
				return err
			}
			skills = append(skills, skill)
		}

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
			ID:         uuid.New(),
			UserID:     userID,
			FullName:   dto.FullName,
			Email:      dto.Email,
			Phone:      dto.Phone,
			Location:   dto.Location,
			Experience: experience,
			Education:  education,
		}

		if err := txRepo.Create(resume); err != nil {
			s.log.Error("Failed to save resume", zap.Error(err))
			return err
		}

		if len(skills) > 0 {
			if err := txRepo.AssociateSkills(resume, skills); err != nil {
				s.log.Error("Failed to associate skills", zap.Error(err))
				return err
			}
		}

		file := &models.ResumeFile{
			ResumeID: resume.ID,
			Path:     path,
			MimeType: "application/pdf",
		}
		if err := txRepo.CreateFile(file); err != nil {
			s.log.Error("Failed to save resume file", zap.Error(err))
			return err
		}

		dto.ID = resume.ID.String()
		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	// Уже вне транзакции — получаем URL
	fileUrl, err := s.GetResumeFileURL(uuid.MustParse(dto.ID))
	if err != nil {
		s.log.Error("Failed to get resume file URL", zap.Error(err))
		return nil, err
	}
	dto.FileURL = fileUrl

	return &dto, nil
}

func (s *ResumeService) GetListResume(userID uuid.UUID) (*response.ResumeListDTO, error) {
	resumes, err := s.repo.GetListRes(userID)
	if err != nil {
		s.log.Error("Failed to get list of resumes", zap.Error(err))
		return nil, err
	}

	var dtos []*response.ResumeListItemDTO
	for _, resume := range *resumes {
		fileUrl, err := s.GetResumeFileURL(resume.ID)
		if err != nil {
			s.log.Error("Failed to get resume file URL", zap.Error(err))
			return nil, err
		}
		dto := &response.ResumeListItemDTO{
			ID:        resume.ID.String(),
			FullName:  resume.FullName,
			FileURL:   fileUrl,
			CreatedAt: resume.CreatedAt,
		}
		dtos = append(dtos, dto)
	}
	return &response.ResumeListDTO{Resumes: dtos}, nil
}

func (s *ResumeService) GetResumeFileURL(resumeID uuid.UUID) (string, error) {
	url, err := s.repo.GetResFileURL(resumeID)
	if err != nil {
		s.log.Error("Failed to get resume file URL", zap.Error(err))
		return "", err
	}
	url = strings.ReplaceAll(url, "./", s.cfg.BaseURL+"/")
	return url, nil
}

func (s *ResumeService) GetResumeByID(userID, resumeID uuid.UUID) (*response.ParsedResumeDTO, error) {
	resume, err := s.repo.GetResumeByID(userID, resumeID)
	if err != nil {
		s.log.Error("Failed to get resume by ID", zap.Error(err))
		return nil, err
	}

	fileUrl, err := s.GetResumeFileURL(resume.ID)
	if err != nil {
		s.log.Error("Failed to get resume file URL", zap.Error(err))
		return nil, err
	}

	var dto response.ParsedResumeDTO
	dto.ID = resume.ID.String()
	dto.FullName = resume.FullName
	dto.Email = resume.Email
	dto.Phone = resume.Phone
	dto.Location = resume.Location
	dto.FileURL = fileUrl
	for _, skill := range resume.Skills {
		dto.Skills = append(dto.Skills, skill.Name)
	}
	for _, exp := range resume.Experience {
		dto.Experience = append(dto.Experience, response.ExperienceDTO{
			Company:     exp.Company,
			Position:    exp.Position,
			StartDate:   exp.StartDate,
			EndDate:     exp.EndDate,
			Description: exp.Description,
		})
	}
	for _, edu := range resume.Education {
		dto.Education = append(dto.Education, response.EducationDTO{
			Institution: edu.Institution,
			Degree:      edu.Degree,
			Field:       edu.Field,
			StartDate:   edu.StartDate,
			EndDate:     edu.EndDate,
		})
	}

	return &dto, nil
}

func (s *ResumeService) DeleteResume(userID, resumeID uuid.UUID) error {
	txErr := s.repo.DB().Transaction(func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)
		_, err := txRepo.GetResumeByID(userID, resumeID)
		if err != nil {
			s.log.Error("Failed to get resume by ID", zap.Error(err))
			return err
		}
		skills, err := txRepo.GetSkillsByResumeID(resumeID)
		if err != nil {
			s.log.Error("Failed to get skills by resume ID", zap.Error(err))
			return err
		}
		for _, skill := range skills {
			if err := txRepo.DeleteSkillFromResume(resumeID, skill.ID); err != nil {
				s.log.Error("Failed to delete skill from resume", zap.Error(err))
				return err
			}
			if err := txRepo.DeleteUnusedSkill(skill.ID); err != nil {
				s.log.Error("Failed to delete unused skill", zap.Error(err))
				return err
			}
		}
		if err := txRepo.DeleteUnusedEdAndEx(resumeID); err != nil {
			s.log.Error("Failed to delete unused education and experience", zap.Error(err))
			return err
		}
		if err := txRepo.DeleteUnusedMatching(resumeID); err != nil {
			s.log.Error("Failed to delete unused matching", zap.Error(err))
			return err
		}
		path, err := txRepo.GetResFileURL(resumeID)
		if err != nil {
			s.log.Error("Failed to get resume file path", zap.Error(err))
			return err
		}
		if err := txRepo.DeleteResumeFile(resumeID); err != nil {
			s.log.Error("Failed to delete resume file", zap.Error(err))
			return err
		}
		os.Remove(path)
		if err := txRepo.DeleteResume(resumeID); err != nil {
			s.log.Error("Failed to delete resume", zap.Error(err))
			return err
		}
		return nil
	})
	if txErr != nil {
		return txErr
	}
	return nil

}
