package service

import (
	"CVMatch/internal/config"
	"CVMatch/internal/jwt"
	"CVMatch/internal/models"
	"CVMatch/internal/repository"
	"errors"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserExists = errors.New("user already exists")

type UserService struct {
	repo repository.UserRepositoryI
	log  *zap.Logger
	cfg  *config.Config
}

func NewUserService(repo repository.UserRepositoryI, log *zap.Logger, cfg *config.Config) *UserService {
	return &UserService{
		repo: repo,
		log:  log,
		cfg:  cfg,
	}
}

func (s *UserService) Register(name, email, password string) (*models.User, error) {
	if _, err := s.repo.FindByEmail(email); err == nil {
		s.log.Warn("User already exists", zap.String("email", email))
		return nil, ErrUserExists
	}

	hashPassword, err := hashedPassword(password)
	if err != nil {
		s.log.Error("Failed to hash password", zap.Error(err))
		return nil, err
	}

	user := &models.User{
		Nickname: name,
		Email:    email,
		Password: hashPassword,
	}

	if err := s.repo.Create(user); err != nil {
		s.log.Error("Failed to create user", zap.Error(err))
		return nil, err
	}

	return user, nil
}

func hashedPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

var ErrUserNotFound = errors.New("user not found")
var ErrInvalidPassword = errors.New("invalid password")

func (s *UserService) Login(email, password string) (access, refresh string, err error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		s.log.Warn("User not found", zap.String("email", email))
		return "", "", ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		s.log.Warn("Invalid password", zap.String("email", email))
		return "", "", ErrInvalidPassword
	}

	access, err = jwt.GenerateAccessToken(user.ID.String(), &s.cfg.JWT)
	if err != nil {
		s.log.Error("Failed to generate access token", zap.Error(err))
		return "", "", err
	}
	refresh, err = jwt.GenerateRefreshToken(user.ID.String(), &s.cfg.JWT)
	if err != nil {
		s.log.Error("Failed to generate refresh token", zap.Error(err))
		return "", "", err
	}

	return access, refresh, nil
}

var ErrInvalidToken = errors.New("invalid token")

func (s *UserService) RefreshToken(refreshToken string) (access, refresh string, err error) {
	claims, err := jwt.ParseRefreshToken(refreshToken, s.cfg.JWT.Refresh)
	if err != nil {
		return "", "", ErrInvalidToken
	}

	access, err = jwt.GenerateAccessToken(claims.UserID, &s.cfg.JWT)
	if err != nil {
		return "", "", err
	}

	refresh, err = jwt.GenerateRefreshToken(claims.UserID, &s.cfg.JWT)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *UserService) Profile(userID string) (*models.User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		s.log.Warn("User not found", zap.String("userID", userID))
		return nil, ErrUserNotFound
	}
	return user, nil
}
