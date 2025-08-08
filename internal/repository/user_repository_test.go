package repository

import (
	"CVMatch/internal/models"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Resume{}, &models.Skill{}, &models.ResumeFile{})
	return db
}
func TestUserRepository_CreateAndFind(t *testing.T) {
	db := setupTestDB()
	repo := NewUserRepository(db)

	user := &models.User{
		Email:    "test@example.com",
		Password: "password",
	}

	err := repo.Create(user)
	require.NoError(t, err)

	found, err := repo.FindByEmail("test@example.com")
	require.NoError(t, err)
	require.Equal(t, user.Email, found.Email)

	foundByID, err := repo.FindByID(user.ID.String())
	require.NoError(t, err)
	require.Equal(t, user.ID, foundByID.ID)
}
