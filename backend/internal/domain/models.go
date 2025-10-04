package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserRole defines user roles in the system
type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleManager  UserRole = "manager"
	RoleEmployee UserRole = "employee"
)

// User represents a user in the system
type User struct {
	ID        primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Email     string              `json:"email" bson:"email"`
	Password  string              `json:"-" bson:"password"` // Never expose in JSON
	FirstName string              `json:"first_name" bson:"first_name"`
	LastName  string              `json:"last_name" bson:"last_name"`
	Role      UserRole            `json:"role" bson:"role"`
	CompanyID primitive.ObjectID  `json:"company_id" bson:"company_id"`
	ManagerID *primitive.ObjectID `json:"manager_id,omitempty" bson:"manager_id,omitempty"` // For employees
	IsActive  bool                `json:"is_active" bson:"is_active"`
	CreatedAt time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time           `json:"updated_at" bson:"updated_at"`
}

// Company represents a company/organization
type Company struct {
	ID             primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Name           string              `json:"name" bson:"name"`
	BaseCurrency   string              `json:"base_currency" bson:"base_currency"` // ISO 4217 currency code
	Country        string              `json:"country" bson:"country"`
	AdminUserID    primitive.ObjectID  `json:"admin_user_id" bson:"admin_user_id"`
	ApprovalRuleID *primitive.ObjectID `json:"approval_rule_id,omitempty" bson:"approval_rule_id,omitempty"`
	IsActive       bool                `json:"is_active" bson:"is_active"`
	CreatedAt      time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at" bson:"updated_at"`
}

// ExpenseStatus defines expense statuses
type ExpenseStatus string

const (
	StatusPending  ExpenseStatus = "pending"
	StatusApproved ExpenseStatus = "approved"
	StatusRejected ExpenseStatus = "rejected"
)

// ExpenseCategory defines expense categories
type ExpenseCategory string

const (
	CategoryTravel        ExpenseCategory = "travel"
	CategoryMeals         ExpenseCategory = "meals"
	CategoryAccommodation ExpenseCategory = "accommodation"
	CategoryTransport     ExpenseCategory = "transport"
	CategorySupplies      ExpenseCategory = "supplies"
	CategoryOther         ExpenseCategory = "other"
)

// Expense represents an expense claim
type Expense struct {
	ID                   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID               primitive.ObjectID `json:"user_id" bson:"user_id"`
	CompanyID            primitive.ObjectID `json:"company_id" bson:"company_id"`
	Amount               float64            `json:"amount" bson:"amount"`
	Currency             string             `json:"currency" bson:"currency"`
	ConvertedAmount      float64            `json:"converted_amount" bson:"converted_amount"` // In company's base currency
	ExchangeRate         float64            `json:"exchange_rate" bson:"exchange_rate"`
	Category             ExpenseCategory    `json:"category" bson:"category"`
	Description          string             `json:"description" bson:"description"`
	ExpenseDate          time.Time          `json:"expense_date" bson:"expense_date"`
	ReceiptURL           string             `json:"receipt_url,omitempty" bson:"receipt_url,omitempty"`
	Merchant             string             `json:"merchant,omitempty" bson:"merchant,omitempty"`
	Status               ExpenseStatus      `json:"status" bson:"status"`
	CurrentApprovalLevel int                `json:"current_approval_level" bson:"current_approval_level"`
	CreatedAt            time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt            time.Time          `json:"updated_at" bson:"updated_at"`
}

// ApprovalStatus defines approval statuses
type ApprovalStatus string

const (
	ApprovalPending  ApprovalStatus = "pending"
	ApprovalApproved ApprovalStatus = "approved"
	ApprovalRejected ApprovalStatus = "rejected"
)

// Approval represents an individual approval action
type Approval struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ExpenseID  primitive.ObjectID `json:"expense_id" bson:"expense_id"`
	ApproverID primitive.ObjectID `json:"approver_id" bson:"approver_id"`
	Level      int                `json:"level" bson:"level"` // Approval level in sequence
	Status     ApprovalStatus     `json:"status" bson:"status"`
	Comments   string             `json:"comments,omitempty" bson:"comments,omitempty"`
	ApprovedAt *time.Time         `json:"approved_at,omitempty" bson:"approved_at,omitempty"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}

// ApprovalWithDetails extends Approval with populated expense and user data
// Used for API responses where related data needs to be included
type ApprovalWithDetails struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ExpenseID  primitive.ObjectID `json:"expense_id" bson:"expense_id"`
	ApproverID primitive.ObjectID `json:"approver_id" bson:"approver_id"`
	Level      int                `json:"level" bson:"level"`
	Status     ApprovalStatus     `json:"status" bson:"status"`
	Comments   string             `json:"comments,omitempty" bson:"comments,omitempty"`
	ApprovedAt *time.Time         `json:"approved_at,omitempty" bson:"approved_at,omitempty"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
	Expense    *ExpenseWithUser   `json:"expense,omitempty" bson:"expense,omitempty"`
	Approver   *User              `json:"approver,omitempty" bson:"approver,omitempty"`
}

// ExpenseWithUser extends Expense with populated user data
type ExpenseWithUser struct {
	ID                   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID               primitive.ObjectID `json:"user_id" bson:"user_id"`
	CompanyID            primitive.ObjectID `json:"company_id" bson:"company_id"`
	Amount               float64            `json:"amount" bson:"amount"`
	Currency             string             `json:"currency" bson:"currency"`
	ConvertedAmount      float64            `json:"converted_amount" bson:"converted_amount"`
	ExchangeRate         float64            `json:"exchange_rate" bson:"exchange_rate"`
	Category             ExpenseCategory    `json:"category" bson:"category"`
	Description          string             `json:"description" bson:"description"`
	ExpenseDate          time.Time          `json:"expense_date" bson:"expense_date"`
	ReceiptURL           string             `json:"receipt_url,omitempty" bson:"receipt_url,omitempty"`
	Merchant             string             `json:"merchant,omitempty" bson:"merchant,omitempty"`
	Status               ExpenseStatus      `json:"status" bson:"status"`
	CurrentApprovalLevel int                `json:"current_approval_level" bson:"current_approval_level"`
	CreatedAt            time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt            time.Time          `json:"updated_at" bson:"updated_at"`
	User                 *User              `json:"user,omitempty" bson:"user,omitempty"`
}

// ApprovalRuleType defines types of approval rules
type ApprovalRuleType string

const (
	RuleTypeSequential       ApprovalRuleType = "sequential"        // All approvers in sequence
	RuleTypePercentage       ApprovalRuleType = "percentage"        // X% approval required
	RuleTypeSpecificApprover ApprovalRuleType = "specific_approver" // Specific person approval
	RuleTypeHybrid           ApprovalRuleType = "hybrid"            // Combination of rules
)

// ApprovalRule defines approval workflow rules for a company
type ApprovalRule struct {
	ID                  primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	CompanyID           primitive.ObjectID   `json:"company_id" bson:"company_id"`
	Name                string               `json:"name" bson:"name"`
	Type                ApprovalRuleType     `json:"type" bson:"type"`
	SequentialApprovers []primitive.ObjectID `json:"sequential_approvers,omitempty" bson:"sequential_approvers,omitempty"`
	PercentageRequired  *float64             `json:"percentage_required,omitempty" bson:"percentage_required,omitempty"` // e.g., 60.0 for 60%
	SpecificApproverID  *primitive.ObjectID  `json:"specific_approver_id,omitempty" bson:"specific_approver_id,omitempty"`
	MinimumApprovals    int                  `json:"minimum_approvals" bson:"minimum_approvals"`
	MaximumApprovals    int                  `json:"maximum_approvals" bson:"maximum_approvals"`
	AllowedApprovers    []primitive.ObjectID `json:"allowed_approvers,omitempty" bson:"allowed_approvers,omitempty"`
	AmountThresholds    []AmountThreshold    `json:"amount_thresholds,omitempty" bson:"amount_thresholds,omitempty"`
	IsActive            bool                 `json:"is_active" bson:"is_active"`
	CreatedAt           time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt           time.Time            `json:"updated_at" bson:"updated_at"`
}

// AmountThreshold defines different approval rules based on expense amount
type AmountThreshold struct {
	MinAmount         float64              `json:"min_amount" bson:"min_amount"`
	MaxAmount         float64              `json:"max_amount" bson:"max_amount"`
	RequiredApprovers []primitive.ObjectID `json:"required_approvers" bson:"required_approvers"`
}

// OCRResult stores OCR extraction results
type OCRResult struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
	ReceiptURL  string             `json:"receipt_url" bson:"receipt_url"`
	Amount      *float64           `json:"amount,omitempty" bson:"amount,omitempty"`
	Currency    *string            `json:"currency,omitempty" bson:"currency,omitempty"`
	Merchant    *string            `json:"merchant,omitempty" bson:"merchant,omitempty"`
	Date        *time.Time         `json:"date,omitempty" bson:"date,omitempty"`
	Category    *string            `json:"category,omitempty" bson:"category,omitempty"`
	RawText     string             `json:"raw_text" bson:"raw_text"`
	Confidence  float64            `json:"confidence" bson:"confidence"` // OCR confidence score
	ProcessedAt time.Time          `json:"processed_at" bson:"processed_at"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}
