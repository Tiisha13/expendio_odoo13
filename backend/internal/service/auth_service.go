package service

import (
	"context"
	"fmt"
	"time"

	"expensio-backend/internal/config"
	"expensio-backend/internal/domain"
	"expensio-backend/pkg/cache"
	"expensio-backend/pkg/currency"
	jwtUtil "expensio-backend/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo    domain.UserRepository
	companyRepo domain.CompanyRepository
	cfg         *config.Config
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo domain.UserRepository, companyRepo domain.CompanyRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		companyRepo: companyRepo,
		cfg:         cfg,
	}
}

type SignupRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	CompanyName string `json:"company_name"`
	Country     string `json:"country"` // Country code (e.g., "US", "GB")
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User         *domain.User    `json:"user"`
	Company      *domain.Company `json:"company"`
	AccessToken  string          `json:"access_token"`
	RefreshToken string          `json:"refresh_token"`
}

// Signup creates a new user and company
func (s *AuthService) Signup(ctx context.Context, req *SignupRequest) (*AuthResponse, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.FindByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("user with this email already exists")
	}

	// Get base currency for the country
	baseCurrency := "USD" // Default
	if req.Country != "" {
		currency, err := currency.GetCountryCurrency(req.Country, s.cfg)
		if err == nil {
			baseCurrency = currency
		}
	}

	// Create company first
	company := &domain.Company{
		Name:         req.CompanyName,
		BaseCurrency: baseCurrency,
		Country:      req.Country,
		IsActive:     true,
	}

	if err := s.companyRepo.Create(ctx, company); err != nil {
		return nil, fmt.Errorf("failed to create company: %w", err)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create admin user
	user := &domain.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      domain.RoleAdmin,
		CompanyID: company.ID,
		IsActive:  true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Update company with admin user ID
	company.AdminUserID = user.ID
	if err := s.companyRepo.Update(ctx, company); err != nil {
		return nil, fmt.Errorf("failed to update company: %w", err)
	}

	// Fetch updated company to get the latest data
	updatedCompany, err := s.companyRepo.FindByID(ctx, company.ID.Hex())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated company: %w", err)
	}

	// Generate JWT tokens
	accessToken, refreshToken, err := jwtUtil.GenerateTokenPair(
		user.ID.Hex(),
		user.Email,
		string(user.Role),
		updatedCompany.ID.Hex(),
		s.cfg,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Store session in Redis
	sessionKey := "session:" + user.ID.Hex()
	sessionData := map[string]interface{}{
		"user_id":    user.ID.Hex(),
		"email":      user.Email,
		"role":       user.Role,
		"company_id": updatedCompany.ID.Hex(),
		"created_at": time.Now(),
	}
	_ = cache.Set(sessionKey, sessionData, s.cfg.JWT.RefreshTokenExpiry)

	return &AuthResponse{
		User:         user,
		Company:      updatedCompany,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Login authenticates a user
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("user account is inactive")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Get company
	company, err := s.companyRepo.FindByID(ctx, user.CompanyID.Hex())
	if err != nil {
		return nil, fmt.Errorf("company not found")
	}

	// Check if company is active
	if !company.IsActive {
		return nil, fmt.Errorf("company account is inactive")
	}

	// Generate JWT tokens
	accessToken, refreshToken, err := jwtUtil.GenerateTokenPair(
		user.ID.Hex(),
		user.Email,
		string(user.Role),
		company.ID.Hex(),
		s.cfg,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Store session in Redis
	sessionKey := "session:" + user.ID.Hex()
	sessionData := map[string]interface{}{
		"user_id":    user.ID.Hex(),
		"email":      user.Email,
		"role":       user.Role,
		"company_id": company.ID.Hex(),
		"created_at": time.Now(),
	}
	_ = cache.Set(sessionKey, sessionData, s.cfg.JWT.RefreshTokenExpiry)

	return &AuthResponse{
		User:         user,
		Company:      company,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Logout invalidates user session
func (s *AuthService) Logout(ctx context.Context, userID, tokenID string) error {
	// Remove session from Redis
	sessionKey := "session:" + userID
	if err := cache.Delete(sessionKey); err != nil {
		return fmt.Errorf("failed to remove session: %w", err)
	}

	// Blacklist the token
	blacklistKey := "blacklist:token:" + tokenID
	_ = cache.SetString(blacklistKey, "1", s.cfg.JWT.AccessTokenExpiry)

	return nil
}

// RefreshToken generates new access token using refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	// Validate refresh token
	claims, err := jwtUtil.ValidateToken(refreshToken, s.cfg)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// Check if token type is refresh
	if claims.TokenType != "refresh" {
		return "", fmt.Errorf("invalid token type")
	}

	// Check if session exists
	sessionKey := "session:" + claims.UserID
	exists, _ := cache.Exists(sessionKey)
	if !exists {
		return "", fmt.Errorf("session expired")
	}

	// Generate new access token
	accessToken, err := jwtUtil.GenerateAccessToken(
		claims.UserID,
		claims.Email,
		claims.Role,
		claims.CompanyID,
		s.cfg,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	return accessToken, nil
}
