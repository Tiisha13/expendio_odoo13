package handler

import (
	"expensio-backend/internal/config"
	"expensio-backend/internal/service"
	"expensio-backend/pkg/response"
	"expensio-backend/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *service.AuthService
	cfg         *config.Config
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		cfg:         cfg,
	}
}

// Signup handles user signup
// @route POST /api/v1/auth/signup
func (h *AuthHandler) Signup(c *fiber.Ctx) error {
	var req service.SignupRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	// Validate request
	if err := validator.ValidateEmail(req.Email); err != nil {
		return response.ValidationError(c, err.Error())
	}
	if err := validator.ValidatePassword(req.Password); err != nil {
		return response.ValidationError(c, err.Error())
	}
	if err := validator.ValidateName(req.FirstName, "First name"); err != nil {
		return response.ValidationError(c, err.Error())
	}
	if err := validator.ValidateName(req.LastName, "Last name"); err != nil {
		return response.ValidationError(c, err.Error())
	}
	if err := validator.ValidateName(req.CompanyName, "Company name"); err != nil {
		return response.ValidationError(c, err.Error())
	}

	// Create user and company
	authResp, err := h.authService.Signup(c.Context(), &req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, "Signup successful", authResp)
}

// Login handles user login
// @route POST /api/v1/auth/login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req service.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	// Validate request
	if err := validator.ValidateEmail(req.Email); err != nil {
		return response.ValidationError(c, err.Error())
	}
	if req.Password == "" {
		return response.ValidationError(c, "Password is required")
	}

	// Authenticate user
	authResp, err := h.authService.Login(c.Context(), &req)
	if err != nil {
		return response.Unauthorized(c, err.Error())
	}

	return response.OK(c, "Login successful", authResp)
}

// Logout handles user logout
// @route POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	tokenID := c.Locals("tokenID").(string)

	if err := h.authService.Logout(c.Context(), userID, tokenID); err != nil {
		return response.InternalServerError(c, "Failed to logout")
	}

	return response.OK(c, "Logout successful", nil)
}

// RefreshToken handles token refresh
// @route POST /api/v1/auth/refresh
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if req.RefreshToken == "" {
		return response.ValidationError(c, "Refresh token is required")
	}

	// Generate new access token
	accessToken, err := h.authService.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		return response.Unauthorized(c, err.Error())
	}

	return response.OK(c, "Token refreshed successfully", fiber.Map{
		"access_token": accessToken,
	})
}
