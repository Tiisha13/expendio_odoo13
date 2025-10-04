package service

import (
	"context"
	"fmt"
	"time"

	"expensio-backend/internal/config"
	"expensio-backend/internal/domain"
	"expensio-backend/pkg/cache"
	"expensio-backend/pkg/currency"
)

type ExpenseService struct {
	expenseRepo     domain.ExpenseRepository
	userRepo        domain.UserRepository
	companyRepo     domain.CompanyRepository
	approvalService *ApprovalService
	cfg             *config.Config
}

// NewExpenseService creates a new expense service
func NewExpenseService(
	expenseRepo domain.ExpenseRepository,
	userRepo domain.UserRepository,
	companyRepo domain.CompanyRepository,
	cfg *config.Config,
) *ExpenseService {
	return &ExpenseService{
		expenseRepo: expenseRepo,
		userRepo:    userRepo,
		companyRepo: companyRepo,
		cfg:         cfg,
	}
}

// SetApprovalService sets the approval service (to avoid circular dependency)
func (s *ExpenseService) SetApprovalService(approvalService *ApprovalService) {
	s.approvalService = approvalService
}

type CreateExpenseRequest struct {
	Amount      float64                `json:"amount"`
	Currency    string                 `json:"currency"`
	Category    domain.ExpenseCategory `json:"category"`
	Description string                 `json:"description"`
	ExpenseDate time.Time              `json:"expense_date"`
	ReceiptURL  string                 `json:"receipt_url,omitempty"`
	Merchant    string                 `json:"merchant,omitempty"`
}

// CreateExpense creates a new expense with currency conversion
func (s *ExpenseService) CreateExpense(ctx context.Context, userID string, req *CreateExpenseRequest) (*domain.Expense, error) {
	// Get user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Get company
	company, err := s.companyRepo.FindByID(ctx, user.CompanyID.Hex())
	if err != nil {
		return nil, fmt.Errorf("company not found")
	}

	// Convert currency to company's base currency
	convertedAmount, exchangeRate, err := currency.ConvertCurrency(
		req.Amount,
		req.Currency,
		company.BaseCurrency,
		s.cfg,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to convert currency: %w", err)
	}

	// Create expense
	expense := &domain.Expense{
		UserID:               user.ID,
		CompanyID:            company.ID,
		Amount:               req.Amount,
		Currency:             req.Currency,
		ConvertedAmount:      convertedAmount,
		ExchangeRate:         exchangeRate,
		Category:             req.Category,
		Description:          req.Description,
		ExpenseDate:          req.ExpenseDate,
		ReceiptURL:           req.ReceiptURL,
		Merchant:             req.Merchant,
		Status:               domain.StatusPending,
		CurrentApprovalLevel: 0,
	}

	if err := s.expenseRepo.Create(ctx, expense); err != nil {
		return nil, fmt.Errorf("failed to create expense: %w", err)
	}

	fmt.Printf("💰 Expense created: %s (Status: %s)\n", expense.ID.Hex(), expense.Status)

	// Initialize approval workflow
	if s.approvalService != nil {
		fmt.Printf("🔄 Approval service available, initializing approvals...\n")
		if err := s.approvalService.InitializeApprovals(ctx, expense); err != nil {
			// Log error but don't fail expense creation
			fmt.Printf("⚠️  Warning: Failed to initialize approvals for expense %s: %v\n", expense.ID.Hex(), err)
		}
	} else {
		fmt.Printf("❌ CRITICAL: Approval service is NULL! Approvals will NOT be created!\n")
		fmt.Printf("   This means SetApprovalService was never called in routes.go\n")
		fmt.Printf("   Backend needs to be restarted!\n")
	}

	// Invalidate relevant caches
	s.invalidateExpenseCaches(user.CompanyID.Hex(), userID)

	return expense, nil
}

// GetExpenseByID retrieves an expense by ID
func (s *ExpenseService) GetExpenseByID(ctx context.Context, expenseID string) (*domain.Expense, error) {
	expense, err := s.expenseRepo.FindByID(ctx, expenseID)
	if err != nil {
		return nil, fmt.Errorf("expense not found")
	}
	return expense, nil
}

// GetUserExpenses retrieves expenses for a user with pagination and caching
func (s *ExpenseService) GetUserExpenses(ctx context.Context, userID string, page, limit int) ([]*domain.Expense, int64, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("expenses:user:%s:page:%d:limit:%d", userID, page, limit)
	var cachedData struct {
		Expenses []*domain.Expense `json:"expenses"`
		Total    int64             `json:"total"`
	}
	err := cache.Get(cacheKey, &cachedData)
	if err == nil {
		return cachedData.Expenses, cachedData.Total, nil
	}

	// Fetch from database
	expenses, total, err := s.expenseRepo.FindByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch expenses: %w", err)
	}

	// Cache the result
	cachedData.Expenses = expenses
	cachedData.Total = total
	_ = cache.Set(cacheKey, cachedData, s.cfg.Cache.ExpenseListTTL)

	return expenses, total, nil
}

// GetCompanyExpenses retrieves all expenses for a company with pagination and caching
func (s *ExpenseService) GetCompanyExpenses(ctx context.Context, companyID string, page, limit int) ([]*domain.Expense, int64, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("expenses:company:%s:page:%d:limit:%d", companyID, page, limit)
	var cachedData struct {
		Expenses []*domain.Expense `json:"expenses"`
		Total    int64             `json:"total"`
	}
	err := cache.Get(cacheKey, &cachedData)
	if err == nil {
		return cachedData.Expenses, cachedData.Total, nil
	}

	// Fetch from database
	expenses, total, err := s.expenseRepo.FindByCompanyID(ctx, companyID, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch expenses: %w", err)
	}

	// Cache the result
	cachedData.Expenses = expenses
	cachedData.Total = total
	_ = cache.Set(cacheKey, cachedData, s.cfg.Cache.ExpenseListTTL)

	return expenses, total, nil
}

// UpdateExpense updates an expense (before approval)
func (s *ExpenseService) UpdateExpense(ctx context.Context, expenseID string, req *CreateExpenseRequest) error {
	// Get existing expense
	expense, err := s.expenseRepo.FindByID(ctx, expenseID)
	if err != nil {
		return fmt.Errorf("expense not found")
	}

	// Only allow updates if status is pending
	if expense.Status != domain.StatusPending {
		return fmt.Errorf("cannot update expense that is already %s", expense.Status)
	}

	// Get company for currency conversion
	company, err := s.companyRepo.FindByID(ctx, expense.CompanyID.Hex())
	if err != nil {
		return fmt.Errorf("company not found")
	}

	// Convert currency
	convertedAmount, exchangeRate, err := currency.ConvertCurrency(
		req.Amount,
		req.Currency,
		company.BaseCurrency,
		s.cfg,
	)
	if err != nil {
		return fmt.Errorf("failed to convert currency: %w", err)
	}

	// Update expense fields
	expense.Amount = req.Amount
	expense.Currency = req.Currency
	expense.ConvertedAmount = convertedAmount
	expense.ExchangeRate = exchangeRate
	expense.Category = req.Category
	expense.Description = req.Description
	expense.ExpenseDate = req.ExpenseDate
	expense.ReceiptURL = req.ReceiptURL
	expense.Merchant = req.Merchant

	if err := s.expenseRepo.Update(ctx, expense); err != nil {
		return fmt.Errorf("failed to update expense: %w", err)
	}

	// Invalidate caches
	s.invalidateExpenseCaches(expense.CompanyID.Hex(), expense.UserID.Hex())

	return nil
}

// DeleteExpense deletes an expense (before approval)
func (s *ExpenseService) DeleteExpense(ctx context.Context, expenseID string) error {
	// Get existing expense
	expense, err := s.expenseRepo.FindByID(ctx, expenseID)
	if err != nil {
		return fmt.Errorf("expense not found")
	}

	// Only allow deletion if status is pending
	if expense.Status != domain.StatusPending {
		return fmt.Errorf("cannot delete expense that is already %s", expense.Status)
	}

	if err := s.expenseRepo.Delete(ctx, expenseID); err != nil {
		return fmt.Errorf("failed to delete expense: %w", err)
	}

	// Invalidate caches
	s.invalidateExpenseCaches(expense.CompanyID.Hex(), expense.UserID.Hex())

	return nil
}

// GetPendingExpenses retrieves pending expenses for a company with caching
func (s *ExpenseService) GetPendingExpenses(ctx context.Context, companyID string) ([]*domain.Expense, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("expenses:pending:company:%s", companyID)
	var cachedExpenses []*domain.Expense
	err := cache.Get(cacheKey, &cachedExpenses)
	if err == nil {
		return cachedExpenses, nil
	}

	// Fetch from database
	expenses, err := s.expenseRepo.FindPendingByCompanyID(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pending expenses: %w", err)
	}

	// Cache the result
	_ = cache.Set(cacheKey, expenses, s.cfg.Cache.PendingApprovalsTTL)

	return expenses, nil
}

// invalidateExpenseCaches invalidates all expense-related caches
func (s *ExpenseService) invalidateExpenseCaches(companyID, userID string) {
	// Invalidate user expense caches
	_ = cache.DeletePattern(fmt.Sprintf("expenses:user:%s:*", userID))

	// Invalidate company expense caches
	_ = cache.DeletePattern(fmt.Sprintf("expenses:company:%s:*", companyID))

	// Invalidate pending expenses cache
	pendingKey := fmt.Sprintf("expenses:pending:company:%s", companyID)
	_ = cache.Delete(pendingKey)
}

// CreateExpenseFromOCR creates an expense from OCR results
func (s *ExpenseService) CreateExpenseFromOCR(ctx context.Context, userID string, ocrResult *domain.OCRResult) (*domain.Expense, error) {
	// Build expense request from OCR data
	req := &CreateExpenseRequest{
		Amount:      *ocrResult.Amount,
		Currency:    *ocrResult.Currency,
		Category:    domain.ExpenseCategory(*ocrResult.Category),
		Description: fmt.Sprintf("Auto-extracted from receipt: %s", *ocrResult.Merchant),
		ExpenseDate: *ocrResult.Date,
		ReceiptURL:  ocrResult.ReceiptURL,
		Merchant:    *ocrResult.Merchant,
	}

	// Use default values if OCR extraction failed
	if req.Amount == 0 {
		return nil, fmt.Errorf("amount is required")
	}
	if req.Currency == "" {
		req.Currency = "USD" // Default currency
	}
	if req.Category == "" {
		req.Category = domain.CategoryOther
	}
	if req.ExpenseDate.IsZero() {
		req.ExpenseDate = time.Now()
	}

	return s.CreateExpense(ctx, userID, req)
}
