package ocr

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"expensio-backend/internal/config"
	"expensio-backend/internal/domain"
	"expensio-backend/pkg/cache"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OCRService handles OCR operations
type OCRService struct {
	cfg *config.Config
}

// NewOCRService creates a new OCR service
func NewOCRService(cfg *config.Config) *OCRService {
	return &OCRService{cfg: cfg}
}

// ProcessReceipt processes a receipt image and extracts expense details
func (s *OCRService) ProcessReceipt(receiptPath string, userID string) (*domain.OCRResult, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("ocr:%s", receiptPath)
	var cachedResult domain.OCRResult
	err := cache.Get(cacheKey, &cachedResult)
	if err == nil {
		return &cachedResult, nil
	}

	// Run Tesseract OCR
	rawText, err := s.runTesseract(receiptPath)
	if err != nil {
		return nil, fmt.Errorf("OCR processing failed: %w", err)
	}

	// Extract structured data from raw text
	result := s.extractData(rawText, receiptPath, userID)

	// Cache the result
	_ = cache.Set(cacheKey, result, s.cfg.Cache.OCRResultTTL)

	return result, nil
}

// runTesseract executes Tesseract OCR on the image
func (s *OCRService) runTesseract(imagePath string) (string, error) {
	// Create output file path
	outputBase := filepath.Join(s.cfg.OCR.TempDir, "ocr_output")

	// Run Tesseract
	cmd := exec.Command(s.cfg.OCR.TesseractPath, imagePath, outputBase)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("tesseract execution failed: %w, output: %s", err, string(output))
	}

	// Read the output file (Tesseract appends .txt)
	outputFile := outputBase + ".txt"
	content, err := exec.Command("cat", outputFile).Output()
	if err != nil {
		return "", fmt.Errorf("failed to read OCR output: %w", err)
	}

	return string(content), nil
}

// extractData extracts structured data from OCR raw text
func (s *OCRService) extractData(rawText, receiptURL, userID string) *domain.OCRResult {
	userObjID, _ := primitive.ObjectIDFromHex(userID)

	result := &domain.OCRResult{
		UserID:      userObjID,
		ReceiptURL:  receiptURL,
		RawText:     rawText,
		Confidence:  0.85, // Placeholder confidence score
		ProcessedAt: time.Now(),
		CreatedAt:   time.Now(),
	}

	// Extract amount
	if amount := extractAmount(rawText); amount != nil {
		result.Amount = amount
	}

	// Extract merchant
	if merchant := extractMerchant(rawText); merchant != nil {
		result.Merchant = merchant
	}

	// Extract date
	if date := extractDate(rawText); date != nil {
		result.Date = date
	}

	// Extract currency
	if currency := extractCurrency(rawText); currency != nil {
		result.Currency = currency
	}

	// Categorize based on merchant/text
	if category := categorizeExpense(rawText); category != nil {
		result.Category = category
	}

	return result
}

// extractAmount extracts monetary amount from text
func extractAmount(text string) *float64 {
	// Common patterns: $123.45, 123.45, €123,45
	patterns := []string{
		`[\$€£¥]?\s*(\d+[,.]?\d*\.?\d+)`,
		`(?i)total:?\s*[\$€£¥]?\s*(\d+[,.]?\d*\.?\d+)`,
		`(?i)amount:?\s*[\$€£¥]?\s*(\d+[,.]?\d*\.?\d+)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(text)
		if len(matches) > 1 {
			amountStr := strings.ReplaceAll(matches[1], ",", "")
			if amount, err := strconv.ParseFloat(amountStr, 64); err == nil && amount > 0 {
				return &amount
			}
		}
	}

	return nil
}

// extractMerchant extracts merchant name from text
func extractMerchant(text string) *string {
	// Typically the first few lines contain merchant name
	lines := strings.Split(text, "\n")
	if len(lines) > 0 {
		merchant := strings.TrimSpace(lines[0])
		if len(merchant) > 3 && len(merchant) < 100 {
			return &merchant
		}
	}
	return nil
}

// extractDate extracts date from text
func extractDate(text string) *time.Time {
	// Common date patterns
	datePatterns := []string{
		`(\d{1,2}[/-]\d{1,2}[/-]\d{2,4})`,
		`(\d{4}[/-]\d{1,2}[/-]\d{1,2})`,
		`(?i)(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)[a-z]*\s+\d{1,2},?\s+\d{4}`,
	}

	for _, pattern := range datePatterns {
		re := regexp.MustCompile(pattern)
		match := re.FindString(text)
		if match != "" {
			// Try to parse the date
			layouts := []string{
				"01/02/2006",
				"2006-01-02",
				"01-02-2006",
				"January 2, 2006",
				"Jan 2, 2006",
			}

			for _, layout := range layouts {
				if parsedDate, err := time.Parse(layout, match); err == nil {
					return &parsedDate
				}
			}
		}
	}

	return nil
}

// extractCurrency extracts currency code from text
func extractCurrency(text string) *string {
	currencySymbols := map[string]string{
		"$": "USD",
		"€": "EUR",
		"£": "GBP",
		"¥": "JPY",
		"₹": "INR",
	}

	for symbol, code := range currencySymbols {
		if strings.Contains(text, symbol) {
			return &code
		}
	}

	// Look for explicit currency codes
	re := regexp.MustCompile(`\b(USD|EUR|GBP|JPY|INR|CAD|AUD|CHF)\b`)
	if match := re.FindString(text); match != "" {
		return &match
	}

	return nil
}

// categorizeExpense categorizes expense based on text content
func categorizeExpense(text string) *string {
	text = strings.ToLower(text)

	categories := map[string][]string{
		"travel":        {"flight", "airline", "airport", "ticket", "travel", "booking"},
		"meals":         {"restaurant", "cafe", "coffee", "food", "dining", "lunch", "dinner", "breakfast"},
		"accommodation": {"hotel", "motel", "accommodation", "lodging", "stay", "airbnb"},
		"transport":     {"taxi", "uber", "lyft", "bus", "train", "metro", "transport", "parking"},
		"supplies":      {"office", "supplies", "stationery", "equipment"},
	}

	for category, keywords := range categories {
		for _, keyword := range keywords {
			if strings.Contains(text, keyword) {
				return &category
			}
		}
	}

	other := "other"
	return &other
}
