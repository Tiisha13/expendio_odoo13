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

type approvalRuleRepository struct {
	collection *mongo.Collection
}

// NewApprovalRuleRepository creates a new approval rule repository
func NewApprovalRuleRepository() domain.ApprovalRuleRepository {
	return &approvalRuleRepository{
		collection: database.GetCollection("approval_rules"),
	}
}

func (r *approvalRuleRepository) Create(ctx context.Context, rule *domain.ApprovalRule) error {
	rule.CreatedAt = time.Now()
	rule.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, rule)
	if err != nil {
		return fmt.Errorf("failed to create approval rule: %w", err)
	}

	rule.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *approvalRuleRepository) FindByID(ctx context.Context, id string) (*domain.ApprovalRule, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid approval rule ID: %w", err)
	}

	var rule domain.ApprovalRule
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&rule)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("approval rule not found")
		}
		return nil, fmt.Errorf("failed to find approval rule: %w", err)
	}

	return &rule, nil
}

func (r *approvalRuleRepository) FindByCompanyID(ctx context.Context, companyID string) (*domain.ApprovalRule, error) {
	objectID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return nil, fmt.Errorf("invalid company ID: %w", err)
	}

	var rule domain.ApprovalRule
	err = r.collection.FindOne(ctx, bson.M{
		"company_id": objectID,
		"is_active":  true,
	}).Decode(&rule)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("approval rule not found")
		}
		return nil, fmt.Errorf("failed to find approval rule: %w", err)
	}

	return &rule, nil
}

func (r *approvalRuleRepository) Update(ctx context.Context, rule *domain.ApprovalRule) error {
	rule.UpdatedAt = time.Now()

	update := bson.M{
		"$set": rule,
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": rule.ID}, update)
	if err != nil {
		return fmt.Errorf("failed to update approval rule: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("approval rule not found")
	}

	return nil
}

func (r *approvalRuleRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid approval rule ID: %w", err)
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete approval rule: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("approval rule not found")
	}

	return nil
}
