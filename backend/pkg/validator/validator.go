package validator

import (
	"fmt"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// ValidateEmail validates email format
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

// ValidatePassword validates password strength
func ValidatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("password is required")
	}
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	return nil
}

// ValidateName validates name fields
func ValidateName(name, fieldName string) error {
	if name == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	if len(name) < 2 {
		return fmt.Errorf("%s must be at least 2 characters long", fieldName)
	}
	if len(name) > 50 {
		return fmt.Errorf("%s must be at most 50 characters long", fieldName)
	}
	return nil
}

// ValidateRole validates user role
func ValidateRole(role string) error {
	validRoles := []string{"admin", "manager", "employee"}
	role = strings.ToLower(role)

	for _, validRole := range validRoles {
		if role == validRole {
			return nil
		}
	}

	return fmt.Errorf("invalid role: must be one of %v", validRoles)
}

// ValidateCurrency validates currency code (ISO 4217)
func ValidateCurrency(currency string) error {
	if currency == "" {
		return fmt.Errorf("currency is required")
	}
	if len(currency) != 3 {
		return fmt.Errorf("currency must be a 3-letter ISO 4217 code")
	}
	return nil
}

// ValidateAmount validates expense amount
func ValidateAmount(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}
	if amount > 1000000000 { // 1 billion limit
		return fmt.Errorf("amount exceeds maximum allowed value")
	}
	return nil
}

// ValidateCategory validates expense category
func ValidateCategory(category string) error {
	validCategories := []string{"travel", "meals", "accommodation", "transport", "supplies", "other"}
	category = strings.ToLower(category)

	for _, validCategory := range validCategories {
		if category == validCategory {
			return nil
		}
	}

	return fmt.Errorf("invalid category: must be one of %v", validCategories)
}

// ValidateDescription validates description length
func ValidateDescription(description string) error {
	if description == "" {
		return fmt.Errorf("description is required")
	}
	if len(description) < 5 {
		return fmt.Errorf("description must be at least 5 characters long")
	}
	if len(description) > 500 {
		return fmt.Errorf("description must be at most 500 characters long")
	}
	return nil
}

// ValidateObjectID validates MongoDB ObjectID format
func ValidateObjectID(id string) error {
	if id == "" {
		return fmt.Errorf("ID is required")
	}
	if len(id) != 24 {
		return fmt.Errorf("invalid ID format")
	}
	// Check if it's a valid hex string
	matched, _ := regexp.MatchString("^[a-fA-F0-9]{24}$", id)
	if !matched {
		return fmt.Errorf("invalid ID format")
	}
	return nil
}

// ValidatePagination validates pagination parameters
func ValidatePagination(page, limit int) error {
	if page < 1 {
		return fmt.Errorf("page must be greater than 0")
	}
	if limit < 1 || limit > 100 {
		return fmt.Errorf("limit must be between 1 and 100")
	}
	return nil
}
