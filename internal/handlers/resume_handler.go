package handlers

import (
	"CVMatch/internal/response"
	"CVMatch/internal/service"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ResumeHandler struct {
	service *service.ResumeService
}

func NewResumeHandler(service *service.ResumeService) *ResumeHandler {
	return &ResumeHandler{
		service: service,
	}
}

// UploadResumeHandler godoc
// @Summary Загрузка резюме
// @Description Загрузка резюме для пользователя
// @Security BearerAuth
// @Tags resumes
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Резюме"
// @Success 200 {object} response.ParsedResumeDTO "Успешная загрузка резюме"
// @Failure 400 {object} response.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Ошибка сервера"
// @Router /resumes/upload [post]
func (h *ResumeHandler) UploadResumeHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: "Unauthorized"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	filename := "resume_" + uuid.New().String() + ".pdf"
	path := "./uploads/" + filename
	if err := c.SaveUploadedFile(file, path); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Error saving file"})
		return
	}

	// Парсим userID в uuid.UUID
	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid user id"})
		return
	}

	resume, err := h.service.CreateResumeWithUser(path, userUUID)
	if err != nil {
		os.Remove(path)
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Error creating resume"})
		return
	}

	c.JSON(http.StatusOK, resume)
}

// ListResumesHandler godoc
// @Summary Получение списка резюме
// @Description Получение списка резюме для пользователя
// @Security BearerAuth
// @Tags resumes
// @Produce json
// @Success 200 {object} response.ResumeListDTO "Успешное получение списка резюме"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Ошибка сервера"
// @Router /resumes/list [get]
func (h *ResumeHandler) ListResumesHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: "Unauthorized"})
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid user id"})
		return
	}

	resumes, err := h.service.GetListResume(userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Error getting resume list"})
		return
	}

	c.JSON(http.StatusOK, resumes)
}

// GetResumeHandler godoc
// @Summary Получение резюме по ID
// @Description Получение резюме по ID для пользователя
// @Security BearerAuth
// @Tags resumes
// @Produce json
// @Param id path string true "ID резюме"
// @Success 200 {object} response.ParsedResumeDTO "Успешное получение резюме"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Resume not found"
// @Failure 500 {object} response.ErrorResponse "Ошибка сервера"
// @Router /resumes/{id} [get]
func (h *ResumeHandler) GetResumeHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: "Unauthorized"})
		return
	}

	resumeID := c.Param("id")
	resumeUUID, err := uuid.Parse(resumeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid resume id"})
		return
	}
	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid user id"})
		return
	}

	resume, err := h.service.GetResumeByID(userUUID, resumeUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Resume not found"})
		return
	}

	c.JSON(http.StatusOK, resume)
}

// DeleteResumeHandler godoc
// @Summary Удаление резюме по ID
// @Description Удаление резюме по ID для пользователя
// @Security BearerAuth
// @Tags resumes
// @Produce json
// @Param id path string true "ID резюме"
// @Success 200 {object} response.SuccessResponse "Успешное удаление резюме"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Ошибка сервера"
// @Router /resumes/{id} [delete]
func (h *ResumeHandler) DeleteResumeHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: "Unauthorized"})
		return
	}
	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid user id"})
		return
	}
	resumeID := c.Param("id")
	resumeUUID, err := uuid.Parse(resumeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid resume id"})
		return
	}

	if err := h.service.DeleteResume(userUUID, resumeUUID); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Error deleting resume"})
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{Message: "Resume deleted successfully"})
}
