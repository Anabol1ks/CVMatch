package repository

import (
	"CVMatch/internal/models"
	"testing"

	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupResumeTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.User{}, &models.Resume{}, &models.Skill{}, &models.ResumeFile{}, &models.Experience{}, &models.Education{})
	return db
}

func TestResumeRepository_CreateAndGetResumeByID(t *testing.T) {
	db := setupResumeTestDB()
	repo := NewResumeRepository(db)
	userID := uuid.New()
	resume := &models.Resume{
		ID:        uuid.New(),
		UserID:    userID,
		FullName:  "Test User",
		Email:     "test@example.com",
		Phone:     "+79999999999",
		Location:  "Moscow",
		CreatedAt: time.Now(),
	}
	err := repo.Create(resume)
	require.NoError(t, err)

	got, err := repo.GetResumeByID(userID, resume.ID)
	require.NoError(t, err)
	require.Equal(t, resume.FullName, got.FullName)
	require.Equal(t, resume.Email, got.Email)
}

func TestResumeRepository_GetListRes(t *testing.T) {
	db := setupResumeTestDB()
	repo := NewResumeRepository(db)
	userID := uuid.New()
	resume1 := &models.Resume{ID: uuid.New(), UserID: userID, FullName: "User1", CreatedAt: time.Now()}
	resume2 := &models.Resume{ID: uuid.New(), UserID: userID, FullName: "User2", CreatedAt: time.Now()}
	_ = repo.Create(resume1)
	_ = repo.Create(resume2)

	list, err := repo.GetListRes(userID)
	require.NoError(t, err)
	require.Len(t, *list, 2)
}

func TestResumeRepository_GetResFileURL(t *testing.T) {
	db := setupResumeTestDB()
	repo := NewResumeRepository(db)
	resumeID := uuid.New()
	file := &models.ResumeFile{ID: uuid.New(), ResumeID: resumeID, Path: "./uploads/test.pdf", MimeType: "application/pdf", CreatedAt: time.Now()}
	err := repo.CreateFile(file)
	require.NoError(t, err)

	url, err := repo.GetResFileURL(resumeID)
	require.NoError(t, err)
	require.Equal(t, "./uploads/test.pdf", url)
}

func TestResumeRepository_FirstOrCreateSkill(t *testing.T) {
	db := setupResumeTestDB()
	repo := NewResumeRepository(db)
	skillName := "Go"
	skill, err := repo.FirstOrCreateSkill(skillName)
	require.NoError(t, err)
	require.Equal(t, skillName, skill.Name)

	// Повторный вызов должен вернуть тот же скилл
	skill2, err := repo.FirstOrCreateSkill(skillName)
	require.NoError(t, err)
	require.Equal(t, skill.ID, skill2.ID)
}
