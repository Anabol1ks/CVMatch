package handlers

import (
	"CVMatch/internal/response"
	"CVMatch/internal/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

type UserRegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Register godoc
// @Summary Регистрация пользователя
// @Description Регистрация нового пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param user body UserRegisterRequest true "Параметры регистрации пользователя"
// @Success 201 {object} response.UserRegisterResponse "Успешная регистрация пользователя"
// @Failure 400 {object} response.ErrorResponse "Ошибка валидации"
// @Failure 409 {object} response.ErrorResponse "Пользователь уже существует"
// @Failure 500 {object} response.ErrorResponse "Ошибка сервера"
// @Router /auth/register [post]
func (h *UserHandler) RegisterHandler(c *gin.Context) {
	var req UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	user, err := h.service.Register(req.Name, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserExists):
			c.JSON(http.StatusConflict, response.ErrorResponse{Error: "User already exists"})
		default:
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, response.UserRegisterResponse{
		ID:       user.ID.String(),
		Nickname: user.Nickname,
		Email:    user.Email,
	})
}

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginHandler godoc
// @Summary Вход пользователя
// @Description Вход существующего пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param user body UserLoginRequest true "Параметры входа пользователя"
// @Success 200 {object} response.TokenResponse "Успешный вход пользователя"
// @Failure 400 {object} response.ErrorResponse "Ошибка валидации"
// @Failure 409 {object} response.ErrorResponse "Пользователь уже существует"
// @Failure 500 {object} response.ErrorResponse "Ошибка сервера"
// @Router /auth/login [post]
func (h *UserHandler) LoginHandler(c *gin.Context) {
	var req UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	access, refresh, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			c.JSON(http.StatusNotFound, response.ErrorResponse{Error: "User not found"})
		case errors.Is(err, service.ErrInvalidPassword):
			c.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: "Invalid password"})
		default:
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, response.TokenResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

type UserRefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Refresh godoc
// @Summary Обновление токена
// @Description Обновление refresh-токена
// @Tags users
// @Accept json
// @Produce json
// @Param user body UserRefreshRequest true "Параметры обновления токена"
// @Success 200 {object} response.TokenResponse "Успешное обновление токена"
// @Failure 400 {object} response.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} response.ErrorResponse "Неверный токен"
// @Failure 500 {object} response.ErrorResponse "Ошибка сервера"
// @Router /auth/refresh [post]
func (h *UserHandler) RefreshHandler(c *gin.Context) {
	var req UserRefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	access, refresh, err := h.service.RefreshToken(req.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidToken):
			c.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: "Invalid refresh token"})
		default:
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, response.TokenResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

func (h *UserHandler) ProfileHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: "Unauthorized"})
		return
	}

	user, err := h.service.Profile(userID.(string))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			c.JSON(http.StatusNotFound, response.ErrorResponse{Error: "User not found"})
		default:
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, response.UserProfileResponse{
		ID:       user.ID.String(),
		Nickname: user.Nickname,
		Email:    user.Email,
	})
}
