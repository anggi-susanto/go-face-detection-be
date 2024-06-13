package rest

import (
	"time"

	"github.com/anggi-susanto/go-face-detection-be/domain"
	"github.com/anggi-susanto/go-face-detection-be/internal/queue"
	"github.com/anggi-susanto/go-face-detection-be/internal/storage"
	"github.com/anggi-susanto/go-face-detection-be/photo"
	"github.com/gofiber/fiber/v2"
)

// ResponseError represents an error response
type ResponseError struct {
	Message string `json:"message"`
}

type PhotoHandler interface {
	Upload(c *fiber.Ctx) error
	CheckResult(c *fiber.Ctx) error
	GetPhoto(c *fiber.Ctx) error
}

type photoHandler struct {
	photoService  photo.Service
	photoProducer queue.Producer
}

func NewPhotoHandler(photoService photo.Service, photoProducer queue.Producer) PhotoHandler {
	return &photoHandler{
		photoService:  photoService,
		photoProducer: photoProducer,
	}
}

// Upload handles file upload.
//
// @Summary upload image for face detection
// @Description upload image for face detection
// @Tags Face Detection
// @Accept json
// @Produce json
// @Param photo body domain.Photo true "photp data"
// @Success 201 {object} domain.Photo
// @Failure 400 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /upload [post]
func (h *photoHandler) Upload(c *fiber.Ctx) error {
	_, err := c.FormFile("photo")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid file"})
	}

	filePath, err := storage.SavePhoto(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save photo"})
	}

	photoID := time.Now().UnixNano()

	photo := &domain.Photo{
		ID:            photoID,
		FilePath:      filePath,
		Status:        "pending",
		FacesDetected: 0,
		TimeStamp:     time.Now(),
	}
	if err := c.BodyParser(photo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ResponseError{
			Message: err.Error(),
		})
	}
	if err := h.photoService.Save(c.Context(), photo); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ResponseError{
			Message: err.Error(),
		})
	}
	if err := h.photoProducer.SendToQueue(photoID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ResponseError{
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Photo uploaded successfully",
	})
}

// CheckResult handles photo check result.
//
// @Summary check photo result
// @Description check photo result
// @Tags Face Detection
// @Accept json
// @Produce json
// @Param id path string true "photo id"
// @Success 200 {object} domain.Photo
// @Failure 500 {object} ResponseError
// @Router /result/{id} [get]
func (h *photoHandler) CheckResult(c *fiber.Ctx) error {
	id := c.Params("id")
	photo, err := h.photoService.CheckResult(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ResponseError{
			Message: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(photo)
}

// GetPhoto handles photo get.
//
// @Summary get photo
// @Description get photo
// @Tags Face Detection
// @Accept json
// @Produce json
// @Param id path string true "photo id"
// @Success 200 {object} domain.Photo
// @Failure 500 {object} ResponseError
// @Router /photo/{id} [get]
func (h *photoHandler) GetPhoto(c *fiber.Ctx) error {
	id := c.Params("id")
	photo, err := h.photoService.GetPhoto(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ResponseError{
			Message: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(photo)
}
