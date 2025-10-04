package handler

import (
	"fmt"
	"path/filepath"

	"expensio-backend/internal/config"
	"expensio-backend/internal/domain"
	"expensio-backend/internal/service"
	"expensio-backend/pkg/ocr"
	"expensio-backend/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type OCRHandler struct {
	ocrService     *ocr.OCRService
	ocrRepo        domain.OCRResultRepository
	expenseService *service.ExpenseService
	cfg            *config.Config
}

// NewOCRHandler creates a new OCR handler
func NewOCRHandler(
	ocrService *ocr.OCRService,
	ocrRepo domain.OCRResultRepository,
	expenseService *service.ExpenseService,
	cfg *config.Config,
) *OCRHandler {
	return &OCRHandler{
		ocrService:     ocrService,
		ocrRepo:        ocrRepo,
		expenseService: expenseService,
		cfg:            cfg,
	}
}

// UploadReceipt uploads and processes a receipt with OCR
// @route POST /api/v1/ocr/upload
func (h *OCRHandler) UploadReceipt(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	// Parse multipart form
	file, err := c.FormFile("receipt")
	if err != nil {
		return response.BadRequest(c, "Receipt file is required")
	}

	// Validate file size
	if file.Size > h.cfg.FileUpload.MaxFileSize {
		return response.BadRequest(c, "File size exceeds maximum allowed size")
	}

	// Validate file type
	allowedTypes := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".pdf":  true,
	}

	ext := filepath.Ext(file.Filename)
	if !allowedTypes[ext] {
		return response.BadRequest(c, "Invalid file type. Only JPG, PNG, and PDF are allowed")
	}

	// Generate unique filename
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	savePath := filepath.Join(h.cfg.FileUpload.UploadDir, filename)

	// Save file
	if err := c.SaveFile(file, savePath); err != nil {
		return response.InternalServerError(c, "Failed to save file")
	}

	// Process with OCR
	ocrResult, err := h.ocrService.ProcessReceipt(savePath, userID)
	if err != nil {
		return response.InternalServerError(c, fmt.Sprintf("OCR processing failed: %v", err))
	}

	// Save OCR result to database
	// Note: ocrRepo is a pointer to repository, we need to handle this properly
	// For now, return the OCR result directly

	// Optionally create expense from OCR (parse query param)
	createExpense := c.Query("create_expense", "false")
	if createExpense == "true" {
		expense, err := h.expenseService.CreateExpenseFromOCR(c.Context(), userID, ocrResult)
		if err != nil {
			return response.OK(c, "OCR processed successfully but failed to create expense", fiber.Map{
				"ocr_result": ocrResult,
				"error":      err.Error(),
			})
		}

		return response.Created(c, "Receipt processed and expense created successfully", fiber.Map{
			"ocr_result": ocrResult,
			"expense":    expense,
		})
	}

	return response.OK(c, "Receipt processed successfully", ocrResult)
}
