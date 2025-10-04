package service

import (
	"context"
	"fmt"

	"expensio-backend/internal/config"
	"expensio-backend/internal/domain"
	"expensio-backend/pkg/cache"
)

type ApprovalService struct {
	approvalRepo     domain.ApprovalRepository
	approvalRuleRepo domain.ApprovalRuleRepository
	expenseRepo      domain.ExpenseRepository
	userRepo         domain.UserRepository
	cfg              *config.Config
}

// NewApprovalService creates a new approval service
func NewApprovalService(
	approvalRepo domain.ApprovalRepository,
	approvalRuleRepo domain.ApprovalRuleRepository,
	expenseRepo domain.ExpenseRepository,
	userRepo domain.UserRepository,
	cfg *config.Config,
) *ApprovalService {
	return &ApprovalService{
		approvalRepo:     approvalRepo,
		approvalRuleRepo: approvalRuleRepo,
		expenseRepo:      expenseRepo,
		userRepo:         userRepo,
		cfg:              cfg,
	}
}

type ApprovalActionRequest struct {
	Comments string `json:"comments,omitempty"`
}

// InitializeApprovals creates approval records based on company's approval rule
func (s *ApprovalService) InitializeApprovals(ctx context.Context, expense *domain.Expense) error {
	// Get company's approval rule
	rule, err := s.approvalRuleRepo.FindByCompanyID(ctx, expense.CompanyID.Hex())
	if err != nil {
		// If no rule exists, create a simple approval for user's manager
		return s.createDefaultApproval(ctx, expense)
	}

	switch rule.Type {
	case domain.RuleTypeSequential:
		return s.createSequentialApprovals(ctx, expense, rule)
	case domain.RuleTypePercentage:
		return s.createPercentageApprovals(ctx, expense, rule)
	case domain.RuleTypeSpecificApprover:
		return s.createSpecificApproverApproval(ctx, expense, rule)
	case domain.RuleTypeHybrid:
		return s.createHybridApprovals(ctx, expense, rule)
	default:
		return s.createDefaultApproval(ctx, expense)
	}
}

// createDefaultApproval creates a simple approval for user's manager
func (s *ApprovalService) createDefaultApproval(ctx context.Context, expense *domain.Expense) error {
	user, err := s.userRepo.FindByID(ctx, expense.UserID.Hex())
	if err != nil {
		return err
	}

	if user.ManagerID == nil {
		// If no manager, auto-approve (for admin users)
		expense.Status = domain.StatusApproved
		return s.expenseRepo.Update(ctx, expense)
	}

	approval := &domain.Approval{
		ExpenseID:  expense.ID,
		ApproverID: *user.ManagerID,
		Level:      1,
		Status:     domain.ApprovalPending,
	}

	return s.approvalRepo.Create(ctx, approval)
}

// createSequentialApprovals creates sequential approval workflow
func (s *ApprovalService) createSequentialApprovals(ctx context.Context, expense *domain.Expense, rule *domain.ApprovalRule) error {
	for i, approverID := range rule.SequentialApprovers {
		approval := &domain.Approval{
			ExpenseID:  expense.ID,
			ApproverID: approverID,
			Level:      i + 1,
			Status:     domain.ApprovalPending,
		}
		if err := s.approvalRepo.Create(ctx, approval); err != nil {
			return err
		}
	}
	return nil
}

// createPercentageApprovals creates approvals for percentage-based rule
func (s *ApprovalService) createPercentageApprovals(ctx context.Context, expense *domain.Expense, rule *domain.ApprovalRule) error {
	for i, approverID := range rule.AllowedApprovers {
		approval := &domain.Approval{
			ExpenseID:  expense.ID,
			ApproverID: approverID,
			Level:      i + 1,
			Status:     domain.ApprovalPending,
		}
		if err := s.approvalRepo.Create(ctx, approval); err != nil {
			return err
		}
	}
	return nil
}

// createSpecificApproverApproval creates approval for specific approver
func (s *ApprovalService) createSpecificApproverApproval(ctx context.Context, expense *domain.Expense, rule *domain.ApprovalRule) error {
	if rule.SpecificApproverID == nil {
		return fmt.Errorf("specific approver not configured")
	}

	approval := &domain.Approval{
		ExpenseID:  expense.ID,
		ApproverID: *rule.SpecificApproverID,
		Level:      1,
		Status:     domain.ApprovalPending,
	}

	return s.approvalRepo.Create(ctx, approval)
}

// createHybridApprovals creates approvals for hybrid rule
func (s *ApprovalService) createHybridApprovals(ctx context.Context, expense *domain.Expense, rule *domain.ApprovalRule) error {
	// Hybrid combines sequential and percentage
	// First create sequential approvals
	if err := s.createSequentialApprovals(ctx, expense, rule); err != nil {
		return err
	}

	// Then create percentage-based approvals if configured
	if len(rule.AllowedApprovers) > 0 {
		return s.createPercentageApprovals(ctx, expense, rule)
	}

	return nil
}

// ApproveExpense approves an expense
func (s *ApprovalService) ApproveExpense(ctx context.Context, expenseID, approverID string, req *ApprovalActionRequest) error {
	// Get expense
	expense, err := s.expenseRepo.FindByID(ctx, expenseID)
	if err != nil {
		return fmt.Errorf("expense not found")
	}

	if expense.Status != domain.StatusPending {
		return fmt.Errorf("expense is already %s", expense.Status)
	}

	// Get all approvals for this expense
	approvals, err := s.approvalRepo.FindByExpenseID(ctx, expenseID)
	if err != nil {
		return fmt.Errorf("failed to fetch approvals: %w", err)
	}

	// Find the approval for this approver
	var currentApproval *domain.Approval
	for _, approval := range approvals {
		if approval.ApproverID.Hex() == approverID && approval.Status == domain.ApprovalPending {
			currentApproval = approval
			break
		}
	}

	if currentApproval == nil {
		return fmt.Errorf("no pending approval found for this user")
	}

	// Update approval status
	currentApproval.Status = domain.ApprovalApproved
	currentApproval.Comments = req.Comments
	if err := s.approvalRepo.Update(ctx, currentApproval); err != nil {
		return fmt.Errorf("failed to update approval: %w", err)
	}

	// Check if expense should be auto-approved based on rules
	shouldAutoApprove, err := s.checkAutoApproval(ctx, expense, approvals)
	if err != nil {
		return fmt.Errorf("failed to check auto-approval: %w", err)
	}

	if shouldAutoApprove {
		expense.Status = domain.StatusApproved
		if err := s.expenseRepo.UpdateStatus(ctx, expenseID, domain.StatusApproved); err != nil {
			return fmt.Errorf("failed to approve expense: %w", err)
		}
	} else {
		// Increment approval level
		expense.CurrentApprovalLevel++
		if err := s.expenseRepo.Update(ctx, expense); err != nil {
			return fmt.Errorf("failed to update expense: %w", err)
		}
	}

	// Invalidate caches
	s.invalidateApprovalCaches(expense.CompanyID.Hex(), approverID)

	return nil
}

// RejectExpense rejects an expense
func (s *ApprovalService) RejectExpense(ctx context.Context, expenseID, approverID string, req *ApprovalActionRequest) error {
	// Get expense
	expense, err := s.expenseRepo.FindByID(ctx, expenseID)
	if err != nil {
		return fmt.Errorf("expense not found")
	}

	if expense.Status != domain.StatusPending {
		return fmt.Errorf("expense is already %s", expense.Status)
	}

	// Get all approvals for this expense
	approvals, err := s.approvalRepo.FindByExpenseID(ctx, expenseID)
	if err != nil {
		return fmt.Errorf("failed to fetch approvals: %w", err)
	}

	// Find the approval for this approver
	var currentApproval *domain.Approval
	for _, approval := range approvals {
		if approval.ApproverID.Hex() == approverID && approval.Status == domain.ApprovalPending {
			currentApproval = approval
			break
		}
	}

	if currentApproval == nil {
		return fmt.Errorf("no pending approval found for this user")
	}

	// Update approval status
	currentApproval.Status = domain.ApprovalRejected
	currentApproval.Comments = req.Comments
	if err := s.approvalRepo.Update(ctx, currentApproval); err != nil {
		return fmt.Errorf("failed to update approval: %w", err)
	}

	// Reject the expense (one rejection rejects all)
	if err := s.expenseRepo.UpdateStatus(ctx, expenseID, domain.StatusRejected); err != nil {
		return fmt.Errorf("failed to reject expense: %w", err)
	}

	// Invalidate caches
	s.invalidateApprovalCaches(expense.CompanyID.Hex(), approverID)

	return nil
}

// checkAutoApproval checks if expense should be auto-approved based on rules
func (s *ApprovalService) checkAutoApproval(ctx context.Context, expense *domain.Expense, approvals []*domain.Approval) (bool, error) {
	// Get approval rule
	rule, err := s.approvalRuleRepo.FindByCompanyID(ctx, expense.CompanyID.Hex())
	if err != nil {
		// No rule, check if all approvals are complete
		return s.checkAllApprovalsComplete(approvals), nil
	}

	switch rule.Type {
	case domain.RuleTypeSequential:
		return s.checkSequentialApproval(approvals), nil
	case domain.RuleTypePercentage:
		return s.checkPercentageApproval(approvals, rule), nil
	case domain.RuleTypeSpecificApprover:
		return s.checkSpecificApproverApproval(approvals, rule), nil
	case domain.RuleTypeHybrid:
		return s.checkHybridApproval(approvals, rule), nil
	default:
		return s.checkAllApprovalsComplete(approvals), nil
	}
}

// checkSequentialApproval checks if sequential approvals are complete
func (s *ApprovalService) checkSequentialApproval(approvals []*domain.Approval) bool {
	return s.checkAllApprovalsComplete(approvals)
}

// checkPercentageApproval checks if percentage threshold is met
func (s *ApprovalService) checkPercentageApproval(approvals []*domain.Approval, rule *domain.ApprovalRule) bool {
	if rule.PercentageRequired == nil {
		return false
	}

	approvedCount := 0
	for _, approval := range approvals {
		if approval.Status == domain.ApprovalApproved {
			approvedCount++
		}
	}

	percentage := float64(approvedCount) / float64(len(approvals)) * 100
	return percentage >= *rule.PercentageRequired
}

// checkSpecificApproverApproval checks if specific approver has approved
func (s *ApprovalService) checkSpecificApproverApproval(approvals []*domain.Approval, rule *domain.ApprovalRule) bool {
	if rule.SpecificApproverID == nil {
		return false
	}

	for _, approval := range approvals {
		if approval.ApproverID == *rule.SpecificApproverID && approval.Status == domain.ApprovalApproved {
			return true
		}
	}

	return false
}

// checkHybridApproval checks hybrid rule (combination of rules)
func (s *ApprovalService) checkHybridApproval(approvals []*domain.Approval, rule *domain.ApprovalRule) bool {
	// Check if specific approver approved (highest priority)
	if s.checkSpecificApproverApproval(approvals, rule) {
		return true
	}

	// Check percentage threshold
	if rule.PercentageRequired != nil && s.checkPercentageApproval(approvals, rule) {
		return true
	}

	// Check if all sequential approvals are complete
	return s.checkSequentialApproval(approvals)
}

// checkAllApprovalsComplete checks if all approvals are approved
func (s *ApprovalService) checkAllApprovalsComplete(approvals []*domain.Approval) bool {
	for _, approval := range approvals {
		if approval.Status != domain.ApprovalApproved {
			return false
		}
	}
	return len(approvals) > 0
}

// GetPendingApprovals retrieves pending approvals for an approver with caching
func (s *ApprovalService) GetPendingApprovals(ctx context.Context, approverID string) ([]*domain.Approval, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("approvals:pending:approver:%s", approverID)
	var cachedApprovals []*domain.Approval
	err := cache.Get(cacheKey, &cachedApprovals)
	if err == nil {
		return cachedApprovals, nil
	}

	// Fetch from database
	approvals, err := s.approvalRepo.FindPendingByApproverID(ctx, approverID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pending approvals: %w", err)
	}

	// Cache the result
	_ = cache.Set(cacheKey, approvals, s.cfg.Cache.PendingApprovalsTTL)

	return approvals, nil
}

// GetApprovalHistory retrieves approval history for an expense
func (s *ApprovalService) GetApprovalHistory(ctx context.Context, expenseID string) ([]*domain.Approval, error) {
	approvals, err := s.approvalRepo.FindByExpenseID(ctx, expenseID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch approval history: %w", err)
	}
	return approvals, nil
}

// invalidateApprovalCaches invalidates approval-related caches
func (s *ApprovalService) invalidateApprovalCaches(companyID, approverID string) {
	// Invalidate pending approvals cache
	pendingKey := fmt.Sprintf("approvals:pending:approver:%s", approverID)
	_ = cache.Delete(pendingKey)

	// Invalidate pending expenses cache
	expensesKey := fmt.Sprintf("expenses:pending:company:%s", companyID)
	_ = cache.Delete(expensesKey)

	// Invalidate company expenses cache
	_ = cache.DeletePattern(fmt.Sprintf("expenses:company:%s:*", companyID))
}
