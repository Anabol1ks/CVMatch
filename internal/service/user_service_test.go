package service

import (
	"CVMatch/internal/config"
	"CVMatch/internal/jwt"
	"CVMatch/internal/models"
	"CVMatch/internal/repository/mocks"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestUserService_Register_UserExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryI(ctrl)
	mockRepo.EXPECT().FindByEmail("test@example.com").Return(&models.User{}, nil)

	service := NewUserService(mockRepo, zap.NewNop(), nil)
	_, err := service.Register("name", "test@example.com", "password123")
	require.ErrorIs(t, err, ErrUserExists)
}

func TestUserService_Register_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryI(ctrl)
	mockRepo.EXPECT().FindByEmail("test@example.com").Return(nil, ErrUserNotFound)
	mockRepo.EXPECT().Create(gomock.Any()).Return(nil)

	service := NewUserService(mockRepo, zap.NewNop(), nil)
	user, err := service.Register("name", "test@example.com", "password123")
	require.NoError(t, err)
	require.NotNil(t, user)
}

var cfg = &config.Config{
	JWT: config.JWTConfig{
		Access:     "test",
		AccessExp:  time.Minute * 15,
		Refresh:    "test",
		RefreshExp: time.Hour * 24 * 7,
	},
}

func TestUserService_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hashed, err := hashedPassword("password123")
	require.NoError(t, err)

	mockRepo := mocks.NewMockUserRepositoryI(ctrl)
	mockRepo.EXPECT().FindByEmail("test@example.com").Return(&models.User{
		Email:    "test@example.com",
		Password: string(hashed),
	}, nil)

	service := NewUserService(mockRepo, zap.NewNop(), cfg)
	access, refresh, err := service.Login("test@example.com", "password123")
	require.NoError(t, err)
	require.NotEmpty(t, access)
	require.NotEmpty(t, refresh)
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryI(ctrl)
	mockRepo.EXPECT().FindByEmail("test@example.com").Return(nil, ErrUserNotFound)

	service := NewUserService(mockRepo, zap.NewNop(), cfg)
	_, _, err := service.Login("test@example.com", "password123")
	require.ErrorIs(t, err, ErrUserNotFound)
}

func TestUserService_Login_InvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hashed, err := hashedPassword("correct_password")
	require.NoError(t, err)

	mockRepo := mocks.NewMockUserRepositoryI(ctrl)
	mockRepo.EXPECT().FindByEmail("test@example.com").Return(&models.User{
		Email:    "test@example.com",
		Password: string(hashed),
	}, nil)

	service := NewUserService(mockRepo, zap.NewNop(), cfg)
	_, _, err = service.Login("test@example.com", "wrong_password")
	require.ErrorIs(t, err, ErrInvalidPassword)
}

func TestUserService_RefreshToken_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryI(ctrl)

	userID := uuid.New().String()
	refreshToken, err := jwt.GenerateRefreshToken(userID, &cfg.JWT)
	require.NoError(t, err)

	service := NewUserService(mockRepo, zap.NewNop(), cfg)
	access, refresh, err := service.RefreshToken(refreshToken)
	require.NoError(t, err)
	require.NotEmpty(t, access)
	require.NotEmpty(t, refresh)
}

func TestUserService_RefreshToken_Invalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryI(ctrl)

	service := NewUserService(mockRepo, zap.NewNop(), cfg)
	_, _, err := service.RefreshToken("invalid_refresh_token")
	require.ErrorIs(t, err, ErrInvalidToken)
}

func TestUserService_Profile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryI(ctrl)
	userID := uuid.New()
	mockRepo.EXPECT().FindByID(userID.String()).Return(&models.User{
		ID:       userID,
		Email:    "test@example.com",
		Nickname: "testuser",
	}, nil)

	service := NewUserService(mockRepo, zap.NewNop(), cfg)
	user, err := service.Profile(userID.String())
	require.NoError(t, err)
	require.NotNil(t, user)
}

func TestUserService_Profile_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryI(ctrl)
	userID := uuid.New()
	mockRepo.EXPECT().FindByID(userID.String()).Return(nil, ErrUserNotFound)

	service := NewUserService(mockRepo, zap.NewNop(), cfg)
	user, err := service.Profile(userID.String())
	require.ErrorIs(t, err, ErrUserNotFound)
	require.Nil(t, user)
}
