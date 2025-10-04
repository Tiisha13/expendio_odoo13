package handler

import (
	"expensio-backend/internal/config"
	"expensio-backend/internal/service"
	"expensio-backend/pkg/response"
	"expensio-backend/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type ApprovalHandler struct {
	approvalService *service.ApprovalService
	cfg             *config.Config
}

// NewApprovalHandler creates a new approval handler
func NewApprovalHandler(approvalService *service.ApprovalService, cfg *config.Config) *ApprovalHandler {
	return &ApprovalHandler{
		approvalService: approvalService,
		cfg:             cfg,
	}
}

// GetPendingApprovals retrieves pending approvals for current user
// @route GET /api/v1/approvals/pending
func (h *ApprovalHandler) GetPendingApprovals(c *fiber.Ctx) error {
	approverID := c.Locals("userID").(string)

	approvals, err := h.approvalService.GetPendingApprovals(c.Context(), approverID)
	if err != nil {
		return response.InternalServerError(c, "Failed to fetch pending approvals")
	}

	return response.OK(c, "Pending approvals retrieved successfully", approvals)
}

// ApproveExpense approves an expense
// @route POST /api/v1/approvals/:id/approve
func (h *ApprovalHandler) ApproveExpense(c *fiber.Ctx) error {
	expenseID := c.Params("id")
	approverID := c.Locals("userID").(string)

	if err := validator.ValidateObjectID(expenseID); err != nil {
		return response.BadRequest(c, "Invalid expense ID")
	}

	var req service.ApprovalActionRequest
	if err := c.BodyParser(&req); err != nil {
		// Comments are optional, so we can ignore parse errors
		req.Comments = ""
	}

	if err := h.approvalService.ApproveExpense(c.Context(), expenseID, approverID, &req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, "Expense approved successfully", nil)
}

// RejectExpense rejects an expense
// @route POST /api/v1/approvals/:id/reject
func (h *ApprovalHandler) RejectExpense(c *fiber.Ctx) error {
	expenseID := c.Params("id")
	approverID := c.Locals("userID").(string)

	if err := validator.ValidateObjectID(expenseID); err != nil {
		return response.BadRequest(c, "Invalid expense ID")
	}

	var req service.ApprovalActionRequest
	if err := c.BodyParser(&req); err != nil {
		// Comments are optional, so we can ignore parse errors
		req.Comments = ""
	}

	if err := h.approvalService.RejectExpense(c.Context(), expenseID, approverID, &req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, "Expense rejected successfully", nil)
}

// GetApprovalHistory retrieves approval history for an expense
// @route GET /api/v1/approvals/history/:expenseId
func (h *ApprovalHandler) GetApprovalHistory(c *fiber.Ctx) error {
	expenseID := c.Params("expenseId")

	if err := validator.ValidateObjectID(expenseID); err != nil {
		return response.BadRequest(c, "Invalid expense ID")
	}

	approvals, err := h.approvalService.GetApprovalHistory(c.Context(), expenseID)
	if err != nil {
		return response.InternalServerError(c, "Failed to fetch approval history")
	}

	return response.OK(c, "Approval history retrieved successfully", approvals)
}
