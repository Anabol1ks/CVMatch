package service

import (
	"CVMatch/internal/config"
	"CVMatch/internal/models"
	"CVMatch/internal/repository/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestResumeService_CreateResumeWithUser_ParseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockResumeRepositoryI(ctrl)
	log := zap.NewNop()
	cfg := &config.Config{BaseURL: "http://localhost:8080"}
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	userID := uuid.New()
	fakePath := "test.pdf"

	mockRepo.EXPECT().DB().Return(db).AnyTimes()
	mockParser := mocks.NewMockResumeParserI(ctrl)
	mockParser.EXPECT().ParseResume(gomock.Any(), gomock.Any()).Return("", assert.AnError)

	service := NewResumeService(mockRepo, log, cfg, mockParser)
	dto, err := service.CreateResumeWithUser(fakePath, userID)
	require.Error(t, err)
	require.Nil(t, dto)
}

func TestResumeService_CreateResumeWithUser_UnmarshalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockResumeRepositoryI(ctrl)
	log := zap.NewNop()
	cfg := &config.Config{BaseURL: "http://localhost:8080"}
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	userID := uuid.New()
	fakePath := "test.pdf"

	mockRepo.EXPECT().DB().Return(db).AnyTimes()
	mockParser := mocks.NewMockResumeParserI(ctrl)
	mockParser.EXPECT().ParseResume(gomock.Any(), gomock.Any()).Return("not a json", nil)

	service := NewResumeService(mockRepo, log, cfg, mockParser)
	dto, err := service.CreateResumeWithUser(fakePath, userID)
	require.Error(t, err)
	require.Nil(t, dto)
}

func TestResumeService_CreateResumeWithUser_TransactionError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockResumeRepositoryI(ctrl)
	log := zap.NewNop()
	cfg := &config.Config{BaseURL: "http://localhost:8080"}
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	userID := uuid.New()
	fakePath := "test.pdf"

	mockRepo.EXPECT().DB().Return(db).AnyTimes()
	mockRepo.EXPECT().WithTx(gomock.Any()).Return(mockRepo).AnyTimes()
	mockRepo.EXPECT().FirstOrCreateSkill("Go").Return(&models.Skill{ID: uuid.New(), Name: "Go"}, nil)
	mockRepo.EXPECT().Create(gomock.Any()).Return(assert.AnError)
	mockParser := mocks.NewMockResumeParserI(ctrl)
	mockParser.EXPECT().ParseResume(gomock.Any(), gomock.Any()).Return(`{"full_name":"Иван Иванов","email":"ivan@test.com","phone":"+79999999999","location":"Москва","skills":["Go"],"experience":[],"education":[]}`, nil)

	service := NewResumeService(mockRepo, log, cfg, mockParser)
	dto, err := service.CreateResumeWithUser(fakePath, userID)
	require.Error(t, err)
	require.Nil(t, dto)
}

func TestResumeService_CreateResumeWithUser_GetResFileURLError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockResumeRepositoryI(ctrl)
	log := zap.NewNop()
	cfg := &config.Config{BaseURL: "http://localhost:8080"}
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	userID := uuid.New()
	fakePath := "test.pdf"

	mockRepo.EXPECT().DB().Return(db).AnyTimes()
	mockRepo.EXPECT().WithTx(gomock.Any()).Return(mockRepo).AnyTimes()
	mockRepo.EXPECT().FirstOrCreateSkill("Go").Return(&models.Skill{ID: uuid.New(), Name: "Go"}, nil)
	mockRepo.EXPECT().Create(gomock.Any()).Return(nil)
	mockRepo.EXPECT().AssociateSkills(gomock.Any(), gomock.Any()).Return(nil)
	mockRepo.EXPECT().CreateFile(gomock.Any()).Return(nil)
	mockRepo.EXPECT().GetResFileURL(gomock.Any()).Return("", assert.AnError)
	mockParser := mocks.NewMockResumeParserI(ctrl)
	mockParser.EXPECT().ParseResume(gomock.Any(), gomock.Any()).Return(`{"full_name":"Иван Иванов","email":"ivan@test.com","phone":"+79999999999","location":"Москва","skills":["Go"],"experience":[],"education":[]}`, nil)

	service := NewResumeService(mockRepo, log, cfg, mockParser)
	dto, err := service.CreateResumeWithUser(fakePath, userID)
	require.Error(t, err)
	require.Nil(t, dto)
}
func TestResumeService_CreateResumeWithUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockResumeRepositoryI(ctrl)
	log := zap.NewNop()
	cfg := &config.Config{BaseURL: "http://localhost:8080"}

	// Создаём in-memory SQLite DB для Transaction
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Подготовим входные данные
	userID := uuid.New()
	fakePath := "test.pdf"
	fakeFileURL := "http://localhost:8080/uploads/test.pdf"

	// Мокаем методы репозитория
	mockRepo.EXPECT().DB().Return(db).AnyTimes() // Теперь возвращаем валидный *gorm.DB
	mockRepo.EXPECT().WithTx(gomock.Any()).Return(mockRepo).AnyTimes()
	mockRepo.EXPECT().FirstOrCreateSkill("Go").Return(&models.Skill{ID: uuid.New(), Name: "Go"}, nil)
	mockRepo.EXPECT().Create(gomock.Any()).Return(nil)
	mockRepo.EXPECT().AssociateSkills(gomock.Any(), gomock.Any()).Return(nil)
	mockRepo.EXPECT().CreateFile(gomock.Any()).Return(nil)
	mockRepo.EXPECT().GetResFileURL(gomock.Any()).Return("./uploads/test.pdf", nil)
	mockParser := mocks.NewMockResumeParserI(ctrl)
	mockParser.EXPECT().ParseResume(gomock.Any(), gomock.Any()).Return(`{"full_name":"Иван Иванов","email":"ivan@test.com","phone":"+79999999999","location":"Москва","skills":["Go"],"experience":[],"education":[]}`, nil)

	service := NewResumeService(mockRepo, log, cfg, mockParser)
	dto, err := service.CreateResumeWithUser(fakePath, userID)
	require.NoError(t, err)
	require.NotNil(t, dto)
	require.Equal(t, "Иван Иванов", dto.FullName)
	require.Equal(t, "ivan@test.com", dto.Email)
	require.Equal(t, fakeFileURL, dto.FileURL)
}

func TestResumeService_GetListResume_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockResumeRepositoryI(ctrl)
	log := zap.NewNop()
	cfg := &config.Config{BaseURL: "http://localhost:8080"}

	userID := uuid.New()
	resumes := []models.Resume{{ID: uuid.New(), FullName: "Test User"}}
	mockRepo.EXPECT().GetListRes(userID).Return(&resumes, nil)
	mockRepo.EXPECT().GetResFileURL(resumes[0].ID).Return("./uploads/test.pdf", nil)

	service := NewResumeService(mockRepo, log, cfg, nil)
	dto, err := service.GetListResume(userID)
	require.NoError(t, err)
	require.NotNil(t, dto)
	require.Len(t, dto.Resumes, 1)
	require.Equal(t, "Test User", dto.Resumes[0].FullName)
}

func TestResumeService_GetListResume_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockResumeRepositoryI(ctrl)
	log := zap.NewNop()
	cfg := &config.Config{BaseURL: "http://localhost:8080"}

	userID := uuid.New()
	mockRepo.EXPECT().GetListRes(userID).Return(nil, assert.AnError)

	service := NewResumeService(mockRepo, log, cfg, nil)
	dto, err := service.GetListResume(userID)
	require.Error(t, err)
	require.Nil(t, dto)
}

func TestResumeService_GetResumeFileURL_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockResumeRepositoryI(ctrl)
	log := zap.NewNop()
	cfg := &config.Config{BaseURL: "http://localhost:8080"}

	resumeID := uuid.New()
	mockRepo.EXPECT().GetResFileURL(resumeID).Return("./uploads/test.pdf", nil)

	service := NewResumeService(mockRepo, log, cfg, nil)
	url, err := service.GetResumeFileURL(resumeID)
	require.NoError(t, err)
	require.Equal(t, "http://localhost:8080/uploads/test.pdf", url)
}

func TestResumeService_GetResumeFileURL_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockResumeRepositoryI(ctrl)
	log := zap.NewNop()
	cfg := &config.Config{BaseURL: "http://localhost:8080"}

	resumeID := uuid.New()
	mockRepo.EXPECT().GetResFileURL(resumeID).Return("", assert.AnError)

	service := NewResumeService(mockRepo, log, cfg, nil)
	url, err := service.GetResumeFileURL(resumeID)
	require.Error(t, err)
	require.Empty(t, url)
}

func TestResumeService_GetResumeByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockResumeRepositoryI(ctrl)
	log := zap.NewNop()
	cfg := &config.Config{BaseURL: "http://localhost:8080"}

	userID := uuid.New()
	resumeID := uuid.New()
	resume := &models.Resume{
		ID:       resumeID,
		FullName: "Test User",
		Email:    "test@mail.com",
		Phone:    "+79999999999",
		Location: "Moscow",
		Skills:   []models.Skill{{Name: "Go"}},
	}
	mockRepo.EXPECT().GetResumeByID(userID, resumeID).Return(resume, nil)
	mockRepo.EXPECT().GetResFileURL(resumeID).Return("./uploads/test.pdf", nil)

	service := NewResumeService(mockRepo, log, cfg, nil)
	dto, err := service.GetResumeByID(userID, resumeID)
	require.NoError(t, err)
	require.NotNil(t, dto)
	require.Equal(t, "Test User", dto.FullName)
	require.Equal(t, "Go", dto.Skills[0])
}

func TestResumeService_GetResumeByID_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockResumeRepositoryI(ctrl)
	log := zap.NewNop()
	cfg := &config.Config{BaseURL: "http://localhost:8080"}

	userID := uuid.New()
	resumeID := uuid.New()
	mockRepo.EXPECT().GetResumeByID(userID, resumeID).Return(nil, assert.AnError)

	service := NewResumeService(mockRepo, log, cfg, nil)
	dto, err := service.GetResumeByID(userID, resumeID)
	require.Error(t, err)
	require.Nil(t, dto)
}

func TestResumeService_GetResumeByID_GetFileURLError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockResumeRepositoryI(ctrl)
	log := zap.NewNop()
	cfg := &config.Config{BaseURL: "http://localhost:8080"}

	userID := uuid.New()
	resumeID := uuid.New()
	resume := &models.Resume{ID: resumeID}
	mockRepo.EXPECT().GetResumeByID(userID, resumeID).Return(resume, nil)
	mockRepo.EXPECT().GetResFileURL(resumeID).Return("", assert.AnError)

	service := NewResumeService(mockRepo, log, cfg, nil)
	dto, err := service.GetResumeByID(userID, resumeID)
	require.Error(t, err)
	require.Nil(t, dto)
}

func TestResumeService_DeleteResume_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockResumeRepositoryI(ctrl)
	log := zap.NewNop()
	cfg := &config.Config{BaseURL: "http://localhost:8080"}
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	userID := uuid.New()
	resumeID := uuid.New()
	mockRepo.EXPECT().DB().Return(db).AnyTimes()
	mockRepo.EXPECT().WithTx(gomock.Any()).Return(mockRepo).AnyTimes()
	mockRepo.EXPECT().GetResumeByID(gomock.Any(), gomock.Any()).Return(&models.Resume{}, nil).AnyTimes()
	mockRepo.EXPECT().GetSkillsByResumeID(resumeID).Return([]*models.Skill{}, nil).AnyTimes()
	mockRepo.EXPECT().DeleteSkillFromResume(resumeID, gomock.Any()).Return(nil).AnyTimes()
	mockRepo.EXPECT().DeleteUnusedSkill(gomock.Any()).Return(nil).AnyTimes()
	mockRepo.EXPECT().DeleteUnusedEdAndEx(resumeID).Return(nil).AnyTimes()
	mockRepo.EXPECT().DeleteUnusedMatching(resumeID).Return(nil).AnyTimes()
	mockRepo.EXPECT().GetResFileURL(resumeID).Return("", nil).AnyTimes()
	mockRepo.EXPECT().DeleteResumeFile(resumeID).Return(nil).AnyTimes()
	mockRepo.EXPECT().DeleteResume(resumeID).Return(nil).AnyTimes()

	service := NewResumeService(mockRepo, log, cfg, nil)
	err = service.DeleteResume(userID, resumeID)
	require.NoError(t, err)
}

func TestResumeService_DeleteResume_TransactionError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockResumeRepositoryI(ctrl)
	log := zap.NewNop()
	cfg := &config.Config{BaseURL: "http://localhost:8080"}
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	userID := uuid.New()
	resumeID := uuid.New()
	mockRepo.EXPECT().DB().Return(db).AnyTimes()
	mockRepo.EXPECT().WithTx(gomock.Any()).Return(mockRepo).AnyTimes()
	mockRepo.EXPECT().GetResumeByID(gomock.Any(), gomock.Any()).Return(&models.Resume{}, nil).AnyTimes()
	mockRepo.EXPECT().GetSkillsByResumeID(resumeID).Return([]*models.Skill{}, nil).AnyTimes()
	mockRepo.EXPECT().DeleteUnusedEdAndEx(resumeID).Return(nil).AnyTimes()
	mockRepo.EXPECT().DeleteUnusedMatching(resumeID).Return(nil).AnyTimes()
	mockRepo.EXPECT().GetResFileURL(resumeID).Return("", nil).AnyTimes()
	mockRepo.EXPECT().DeleteResumeFile(resumeID).Return(nil).AnyTimes()
	mockRepo.EXPECT().DeleteResume(resumeID).Return(assert.AnError)

	db.Callback().Delete().Before("gorm:delete").Register("test_error", func(db *gorm.DB) {
		db.AddError(assert.AnError)
	})

	service := NewResumeService(mockRepo, log, cfg, nil)
	err = service.DeleteResume(userID, resumeID)
	require.Error(t, err)
}
