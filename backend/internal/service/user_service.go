package service

import (
	"context"
	"fmt"

	"expensio-backend/internal/config"
	"expensio-backend/internal/domain"
	"expensio-backend/pkg/cache"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo    domain.UserRepository
	companyRepo domain.CompanyRepository
	cfg         *config.Config
}

// NewUserService creates a new user service
func NewUserService(userRepo domain.UserRepository, companyRepo domain.CompanyRepository, cfg *config.Config) *UserService {
	return &UserService{
		userRepo:    userRepo,
		companyRepo: companyRepo,
		cfg:         cfg,
	}
}

type CreateUserRequest struct {
	Email     string          `json:"email"`
	Password  string          `json:"password"`
	FirstName string          `json:"first_name"`
	LastName  string          `json:"last_name"`
	Role      domain.UserRole `json:"role"`
	ManagerID *string         `json:"manager_id,omitempty"`
}

// CreateUser creates a new user (Admin only)
func (s *UserService) CreateUser(ctx context.Context, companyID string, req *CreateUserRequest) (*domain.User, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.FindByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("user with this email already exists")
	}

	// Validate company exists
	companyObjID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return nil, fmt.Errorf("invalid company ID")
	}

	company, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("company not found")
	}

	if !company.IsActive {
		return nil, fmt.Errorf("company is inactive")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &domain.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      req.Role,
		CompanyID: companyObjID,
		IsActive:  true,
	}

	// Set manager if provided
	if req.ManagerID != nil && *req.ManagerID != "" {
		managerObjID, err := primitive.ObjectIDFromHex(*req.ManagerID)
		if err != nil {
			return nil, fmt.Errorf("invalid manager ID")
		}

		// Verify manager exists and is in same company
		manager, err := s.userRepo.FindByID(ctx, *req.ManagerID)
		if err != nil {
			return nil, fmt.Errorf("manager not found")
		}

		if manager.CompanyID != companyObjID {
			return nil, fmt.Errorf("manager must be in the same company")
		}

		user.ManagerID = &managerObjID
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Invalidate user list cache
	cacheKey := fmt.Sprintf("users:company:%s", companyID)
	_ = cache.Delete(cacheKey)

	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// GetUsersByCompanyID retrieves all users in a company with caching
func (s *UserService) GetUsersByCompanyID(ctx context.Context, companyID string) ([]*domain.User, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("users:company:%s", companyID)
	var cachedUsers []*domain.User
	err := cache.Get(cacheKey, &cachedUsers)
	if err == nil {
		return cachedUsers, nil
	}

	// Fetch from database
	users, err := s.userRepo.FindByCompanyID(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	// Cache the result
	_ = cache.Set(cacheKey, users, s.cfg.Cache.DefaultTTL)

	return users, nil
}

// UpdateRole updates a user's role (Admin only)
func (s *UserService) UpdateRole(ctx context.Context, userID string, newRole domain.UserRole) error {
	// Validate user exists
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Update role
	if err := s.userRepo.UpdateRole(ctx, userID, newRole); err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	// Invalidate caches
	cacheKey := fmt.Sprintf("users:company:%s", user.CompanyID.Hex())
	_ = cache.Delete(cacheKey)

	// Invalidate user session to force re-login with new role
	sessionKey := "session:" + userID
	_ = cache.Delete(sessionKey)

	return nil
}

// AssignManager assigns a manager to a user
func (s *UserService) AssignManager(ctx context.Context, userID, managerID string) error {
	// Validate user exists
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Validate manager exists
	manager, err := s.userRepo.FindByID(ctx, managerID)
	if err != nil {
		return fmt.Errorf("manager not found")
	}

	// Verify both are in same company
	if user.CompanyID != manager.CompanyID {
		return fmt.Errorf("user and manager must be in the same company")
	}

	// Assign manager
	if err := s.userRepo.AssignManager(ctx, userID, managerID); err != nil {
		return fmt.Errorf("failed to assign manager: %w", err)
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("users:company:%s", user.CompanyID.Hex())
	_ = cache.Delete(cacheKey)

	return nil
}

// DeleteUser deletes a user (Admin only)
func (s *UserService) DeleteUser(ctx context.Context, userID string) error {
	// Validate user exists
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Delete user
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Invalidate caches
	cacheKey := fmt.Sprintf("users:company:%s", user.CompanyID.Hex())
	_ = cache.Delete(cacheKey)

	sessionKey := "session:" + userID
	_ = cache.Delete(sessionKey)

	return nil
}
