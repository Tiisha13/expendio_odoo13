package domain

import "context"

// UserRepository defines methods for user data access
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByCompanyID(ctx context.Context, companyID string) ([]*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
	UpdateRole(ctx context.Context, id string, role UserRole) error
	AssignManager(ctx context.Context, userID, managerID string) error
}

// CompanyRepository defines methods for company data access
type CompanyRepository interface {
	Create(ctx context.Context, company *Company) error
	FindByID(ctx context.Context, id string) (*Company, error)
	FindByName(ctx context.Context, name string) (*Company, error)
	Update(ctx context.Context, company *Company) error
	Delete(ctx context.Context, id string) error
}

// ExpenseRepository defines methods for expense data access
type ExpenseRepository interface {
	Create(ctx context.Context, expense *Expense) error
	FindByID(ctx context.Context, id string) (*Expense, error)
	FindByUserID(ctx context.Context, userID string, page, limit int) ([]*Expense, int64, error)
	FindByCompanyID(ctx context.Context, companyID string, page, limit int) ([]*Expense, int64, error)
	Update(ctx context.Context, expense *Expense) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status ExpenseStatus) error
	FindPendingByCompanyID(ctx context.Context, companyID string) ([]*Expense, error)
}

// ApprovalRepository defines methods for approval data access
type ApprovalRepository interface {
	Create(ctx context.Context, approval *Approval) error
	FindByID(ctx context.Context, id string) (*Approval, error)
	FindByExpenseID(ctx context.Context, expenseID string) ([]*Approval, error)
	FindPendingByApproverID(ctx context.Context, approverID string) ([]*Approval, error)
	Update(ctx context.Context, approval *Approval) error
	UpdateStatus(ctx context.Context, id string, status ApprovalStatus) error
	CountApprovedByExpenseID(ctx context.Context, expenseID string) (int64, error)
	CountTotalByExpenseID(ctx context.Context, expenseID string) (int64, error)
}

// ApprovalRuleRepository defines methods for approval rule data access
type ApprovalRuleRepository interface {
	Create(ctx context.Context, rule *ApprovalRule) error
	FindByID(ctx context.Context, id string) (*ApprovalRule, error)
	FindByCompanyID(ctx context.Context, companyID string) (*ApprovalRule, error)
	Update(ctx context.Context, rule *ApprovalRule) error
	Delete(ctx context.Context, id string) error
}

// OCRResultRepository defines methods for OCR result data access
type OCRResultRepository interface {
	Create(ctx context.Context, result *OCRResult) error
	FindByID(ctx context.Context, id string) (*OCRResult, error)
	FindByReceiptURL(ctx context.Context, receiptURL string) (*OCRResult, error)
	FindByUserID(ctx context.Context, userID string) ([]*OCRResult, error)
}
