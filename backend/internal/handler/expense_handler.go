package handler

import (
	"strconv"

	"expensio-backend/internal/config"
	"expensio-backend/internal/service"
	"expensio-backend/pkg/response"
	"expensio-backend/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type ExpenseHandler struct {
	expenseService *service.ExpenseService
	cfg            *config.Config
}

// NewExpenseHandler creates a new expense handler
func NewExpenseHandler(expenseService *service.ExpenseService, cfg *config.Config) *ExpenseHandler {
	return &ExpenseHandler{
		expenseService: expenseService,
		cfg:            cfg,
	}
}

// CreateExpense creates a new expense
// @route POST /api/v1/expenses
func (h *ExpenseHandler) CreateExpense(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req service.CreateExpenseRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	// Validate request
	if err := validator.ValidateAmount(req.Amount); err != nil {
		return response.ValidationError(c, err.Error())
	}
	if err := validator.ValidateCurrency(req.Currency); err != nil {
		return response.ValidationError(c, err.Error())
	}
	if err := validator.ValidateCategory(string(req.Category)); err != nil {
		return response.ValidationError(c, err.Error())
	}
	if err := validator.ValidateDescription(req.Description); err != nil {
		return response.ValidationError(c, err.Error())
	}

	// Create expense
	expense, err := h.expenseService.CreateExpense(c.Context(), userID, &req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, "Expense created successfully", expense)
}

// GetExpense retrieves a single expense by ID
// @route GET /api/v1/expenses/:id
func (h *ExpenseHandler) GetExpense(c *fiber.Ctx) error {
	expenseID := c.Params("id")

	if err := validator.ValidateObjectID(expenseID); err != nil {
		return response.BadRequest(c, "Invalid expense ID")
	}

	expense, err := h.expenseService.GetExpenseByID(c.Context(), expenseID)
	if err != nil {
		return response.NotFound(c, "Expense not found")
	}

	return response.OK(c, "Expense retrieved successfully", expense)
}

// GetExpenses retrieves expenses with pagination
// @route GET /api/v1/expenses
func (h *ExpenseHandler) GetExpenses(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	companyID := c.Locals("companyID").(string)
	role := c.Locals("role").(string)

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if err := validator.ValidatePagination(page, limit); err != nil {
		return response.ValidationError(c, err.Error())
	}

	var expenses []*interface{}
	var total int64
	var err error

	// Admins and Managers can see all company expenses
	if role == "admin" || role == "manager" {
		result, t, e := h.expenseService.GetCompanyExpenses(c.Context(), companyID, page, limit)
		expenses = make([]*interface{}, len(result))
		for i, v := range result {
			var temp interface{} = v
			expenses[i] = &temp
		}
		total = t
		err = e
	} else {
		// Employees see only their expenses
		result, t, e := h.expenseService.GetUserExpenses(c.Context(), userID, page, limit)
		expenses = make([]*interface{}, len(result))
		for i, v := range result {
			var temp interface{} = v
			expenses[i] = &temp
		}
		total = t
		err = e
	}

	if err != nil {
		return response.InternalServerError(c, "Failed to fetch expenses")
	}

	meta := fiber.Map{
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": (total + int64(limit) - 1) / int64(limit),
	}

	return response.SuccessWithMeta(c, fiber.StatusOK, "Expenses retrieved successfully", expenses, meta)
}

// UpdateExpense updates an expense
// @route PUT /api/v1/expenses/:id
func (h *ExpenseHandler) UpdateExpense(c *fiber.Ctx) error {
	expenseID := c.Params("id")

	if err := validator.ValidateObjectID(expenseID); err != nil {
		return response.BadRequest(c, "Invalid expense ID")
	}

	var req service.CreateExpenseRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	// Validate request
	if err := validator.ValidateAmount(req.Amount); err != nil {
		return response.ValidationError(c, err.Error())
	}
	if err := validator.ValidateCurrency(req.Currency); err != nil {
		return response.ValidationError(c, err.Error())
	}
	if err := validator.ValidateCategory(string(req.Category)); err != nil {
		return response.ValidationError(c, err.Error())
	}
	if err := validator.ValidateDescription(req.Description); err != nil {
		return response.ValidationError(c, err.Error())
	}

	if err := h.expenseService.UpdateExpense(c.Context(), expenseID, &req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, "Expense updated successfully", nil)
}

// DeleteExpense deletes an expense
// @route DELETE /api/v1/expenses/:id
func (h *ExpenseHandler) DeleteExpense(c *fiber.Ctx) error {
	expenseID := c.Params("id")

	if err := validator.ValidateObjectID(expenseID); err != nil {
		return response.BadRequest(c, "Invalid expense ID")
	}

	if err := h.expenseService.DeleteExpense(c.Context(), expenseID); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, "Expense deleted successfully", nil)
}

// GetPendingExpenses retrieves pending expenses for managers/admins
// @route GET /api/v1/expenses/pending
func (h *ExpenseHandler) GetPendingExpenses(c *fiber.Ctx) error {
	companyID := c.Locals("companyID").(string)

	expenses, err := h.expenseService.GetPendingExpenses(c.Context(), companyID)
	if err != nil {
		return response.InternalServerError(c, "Failed to fetch pending expenses")
	}

	return response.OK(c, "Pending expenses retrieved successfully", expenses)
}
