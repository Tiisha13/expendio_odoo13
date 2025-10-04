package handler

import (
	"expensio-backend/internal/config"
	"expensio-backend/internal/domain"
	"expensio-backend/internal/service"
	"expensio-backend/pkg/response"
	"expensio-backend/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService *service.UserService
	cfg         *config.Config
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *service.UserService, cfg *config.Config) *UserHandler {
	return &UserHandler{
		userService: userService,
		cfg:         cfg,
	}
}

// CreateUser creates a new user (Admin only)
// @route POST /api/v1/users
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	companyID := c.Locals("companyID").(string)

	var req service.CreateUserRequest
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
	if err := validator.ValidateRole(string(req.Role)); err != nil {
		return response.ValidationError(c, err.Error())
	}

	// Create user
	user, err := h.userService.CreateUser(c.Context(), companyID, &req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, "User created successfully", user)
}

// GetUsers retrieves all users in the company
// @route GET /api/v1/users
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	companyID := c.Locals("companyID").(string)

	users, err := h.userService.GetUsersByCompanyID(c.Context(), companyID)
	if err != nil {
		return response.InternalServerError(c, "Failed to fetch users")
	}

	return response.OK(c, "Users retrieved successfully", users)
}

// GetUser retrieves a single user by ID
// @route GET /api/v1/users/:id
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	if err := validator.ValidateObjectID(userID); err != nil {
		return response.BadRequest(c, "Invalid user ID")
	}

	user, err := h.userService.GetUserByID(c.Context(), userID)
	if err != nil {
		return response.NotFound(c, "User not found")
	}

	return response.OK(c, "User retrieved successfully", user)
}

// UpdateUserRole updates a user's role (Admin only)
// @route PUT /api/v1/users/:id/role
func (h *UserHandler) UpdateUserRole(c *fiber.Ctx) error {
	userID := c.Params("id")

	if err := validator.ValidateObjectID(userID); err != nil {
		return response.BadRequest(c, "Invalid user ID")
	}

	var req struct {
		Role domain.UserRole `json:"role"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if err := validator.ValidateRole(string(req.Role)); err != nil {
		return response.ValidationError(c, err.Error())
	}

	if err := h.userService.UpdateRole(c.Context(), userID, req.Role); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, "User role updated successfully", nil)
}

// AssignManager assigns a manager to a user
// @route PUT /api/v1/users/:id/manager
func (h *UserHandler) AssignManager(c *fiber.Ctx) error {
	userID := c.Params("id")

	if err := validator.ValidateObjectID(userID); err != nil {
		return response.BadRequest(c, "Invalid user ID")
	}

	var req struct {
		ManagerID string `json:"manager_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if err := validator.ValidateObjectID(req.ManagerID); err != nil {
		return response.BadRequest(c, "Invalid manager ID")
	}

	if err := h.userService.AssignManager(c.Context(), userID, req.ManagerID); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, "Manager assigned successfully", nil)
}

// DeleteUser deletes a user (Admin only)
// @route DELETE /api/v1/users/:id
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	if err := validator.ValidateObjectID(userID); err != nil {
		return response.BadRequest(c, "Invalid user ID")
	}

	if err := h.userService.DeleteUser(c.Context(), userID); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, "User deleted successfully", nil)
}
