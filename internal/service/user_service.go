package service

import (
	"CVMatch/internal/config"
	"CVMatch/internal/models"
	"CVMatch/internal/repository"
	"errors"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserExists = errors.New("user already exists")

type UserService struct {
	repo *repository.UserRepository
	log  *zap.Logger
	cfg  *config.Config
}

func NewUserService(repo *repository.UserRepository, log *zap.Logger, cfg *config.Config) *UserService {
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
