package repository

import (
	"context"
	"fmt"
	"time"

	"expensio-backend/internal/domain"
	"expensio-backend/pkg/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type approvalRepository struct {
	collection *mongo.Collection
}

// NewApprovalRepository creates a new approval repository
func NewApprovalRepository() domain.ApprovalRepository {
	return &approvalRepository{
		collection: database.GetCollection("approvals"),
	}
}

func (r *approvalRepository) Create(ctx context.Context, approval *domain.Approval) error {
	approval.CreatedAt = time.Now()
	approval.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, approval)
	if err != nil {
		return fmt.Errorf("failed to create approval: %w", err)
	}

	approval.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *approvalRepository) FindByID(ctx context.Context, id string) (*domain.Approval, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid approval ID: %w", err)
	}

	var approval domain.Approval
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&approval)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("approval not found")
		}
		return nil, fmt.Errorf("failed to find approval: %w", err)
	}

	return &approval, nil
}

func (r *approvalRepository) FindByExpenseID(ctx context.Context, expenseID string) ([]*domain.Approval, error) {
	objectID, err := primitive.ObjectIDFromHex(expenseID)
	if err != nil {
		return nil, fmt.Errorf("invalid expense ID: %w", err)
	}

	cursor, err := r.collection.Find(ctx, bson.M{"expense_id": objectID})
	if err != nil {
		return nil, fmt.Errorf("failed to find approvals: %w", err)
	}
	defer cursor.Close(ctx)

	var approvals []*domain.Approval
	if err := cursor.All(ctx, &approvals); err != nil {
		return nil, fmt.Errorf("failed to decode approvals: %w", err)
	}

	return approvals, nil
}

func (r *approvalRepository) FindPendingByApproverID(ctx context.Context, approverID string) ([]*domain.Approval, error) {
	objectID, err := primitive.ObjectIDFromHex(approverID)
	if err != nil {
		return nil, fmt.Errorf("invalid approver ID: %w", err)
	}

	filter := bson.M{
		"approver_id": objectID,
		"status":      domain.ApprovalPending,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find pending approvals: %w", err)
	}
	defer cursor.Close(ctx)

	var approvals []*domain.Approval
	if err := cursor.All(ctx, &approvals); err != nil {
		return nil, fmt.Errorf("failed to decode approvals: %w", err)
	}

	return approvals, nil
}

func (r *approvalRepository) Update(ctx context.Context, approval *domain.Approval) error {
	approval.UpdatedAt = time.Now()

	update := bson.M{
		"$set": approval,
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": approval.ID}, update)
	if err != nil {
		return fmt.Errorf("failed to update approval: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("approval not found")
	}

	return nil
}

func (r *approvalRepository) UpdateStatus(ctx context.Context, id string, status domain.ApprovalStatus) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid approval ID: %w", err)
	}

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"status":      status,
			"approved_at": &now,
			"updated_at":  now,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return fmt.Errorf("failed to update approval status: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("approval not found")
	}

	return nil
}

func (r *approvalRepository) CountApprovedByExpenseID(ctx context.Context, expenseID string) (int64, error) {
	objectID, err := primitive.ObjectIDFromHex(expenseID)
	if err != nil {
		return 0, fmt.Errorf("invalid expense ID: %w", err)
	}

	filter := bson.M{
		"expense_id": objectID,
		"status":     domain.ApprovalApproved,
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count approved approvals: %w", err)
	}

	return count, nil
}

func (r *approvalRepository) CountTotalByExpenseID(ctx context.Context, expenseID string) (int64, error) {
	objectID, err := primitive.ObjectIDFromHex(expenseID)
	if err != nil {
		return 0, fmt.Errorf("invalid expense ID: %w", err)
	}

	filter := bson.M{"expense_id": objectID}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count total approvals: %w", err)
	}

	return count, nil
}
