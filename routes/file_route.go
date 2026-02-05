package routes

import (
	"upload/middleware"
	"upload/services"

	"github.com/gofiber/fiber/v2"
)

type FileHandler struct {
	service *services.FileService
}

func NewFileHandler(service *services.FileService) *FileHandler {
	return &FileHandler{service: service}
}

func (h *FileHandler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api")
	files := api.Group("/files")

	files.Post("/upload", middleware.AuthMiddleware(), h.UploadFile)
	files.Post("/photo", middleware.AuthMiddleware(), h.UploadPhoto)
	files.Post("/certificate", middleware.AuthMiddleware(), h.UploadCertificate)
	files.Get("/", h.GetAllFiles)
	files.Get("/:id", h.GetFileByID)
	files.Delete("/:id", middleware.AuthMiddleware(), h.DeleteFile)
}

func (h *FileHandler) UploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file uploaded",
		})
	}

	userID := c.Locals("userID").(string)

	response, err := h.service.Upload(file, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func (h *FileHandler) UploadPhoto(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file uploaded",
		})
	}

	userID := c.Locals("userID").(string)

	targetUserID := c.FormValue("target_user_id")
	role := c.Locals("role").(string)
	if targetUserID != "" && role == "admin" {
		userID = targetUserID
	}

	response, err := h.service.UploadPhoto(file, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func (h *FileHandler) UploadCertificate(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file uploaded",
		})
	}

	userID := c.Locals("userID").(string)

	targetUserID := c.FormValue("target_user_id")
	role := c.Locals("role").(string)
	if targetUserID != "" && role == "admin" {
		userID = targetUserID
	}

	response, err := h.service.UploadCertificate(file, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func (h *FileHandler) GetAllFiles(c *fiber.Ctx) error {
	responses, err := h.service.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": responses,
	})
}

func (h *FileHandler) GetFileByID(c *fiber.Ctx) error {
	id := c.Params("id")

	response, err := h.service.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(response)
}

func (h *FileHandler) DeleteFile(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.service.Delete(id); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
