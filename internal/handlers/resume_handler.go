package handlers

import (
	"CVMatch/internal/response"
	"CVMatch/internal/service"
	"net/http"

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

	filename := file.Filename + "_" + userID.(string)
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
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Error creating resume"})
		return
	}

	c.JSON(http.StatusOK, resume)
}
