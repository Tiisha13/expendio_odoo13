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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type expenseRepository struct {
	collection *mongo.Collection
}

// NewExpenseRepository creates a new expense repository
func NewExpenseRepository() domain.ExpenseRepository {
	return &expenseRepository{
		collection: database.GetCollection("expenses"),
	}
}

func (r *expenseRepository) Create(ctx context.Context, expense *domain.Expense) error {
	expense.CreatedAt = time.Now()
	expense.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, expense)
	if err != nil {
		return fmt.Errorf("failed to create expense: %w", err)
	}

	expense.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *expenseRepository) FindByID(ctx context.Context, id string) (*domain.Expense, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid expense ID: %w", err)
	}

	var expense domain.Expense
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&expense)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("expense not found")
		}
		return nil, fmt.Errorf("failed to find expense: %w", err)
	}

	return &expense, nil
}

func (r *expenseRepository) FindByUserID(ctx context.Context, userID string, page, limit int) ([]*domain.Expense, int64, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid user ID: %w", err)
	}

	filter := bson.M{"user_id": objectID}

	// Count total documents
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count expenses: %w", err)
	}

	// Calculate skip
	skip := int64((page - 1) * limit)

	// Find with pagination
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(skip).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find expenses: %w", err)
	}
	defer cursor.Close(ctx)

	var expenses []*domain.Expense
	if err := cursor.All(ctx, &expenses); err != nil {
		return nil, 0, fmt.Errorf("failed to decode expenses: %w", err)
	}

	return expenses, total, nil
}

func (r *expenseRepository) FindByCompanyID(ctx context.Context, companyID string, page, limit int) ([]*domain.Expense, int64, error) {
	objectID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid company ID: %w", err)
	}

	filter := bson.M{"company_id": objectID}

	// Count total documents
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count expenses: %w", err)
	}

	// Calculate skip
	skip := int64((page - 1) * limit)

	// Find with pagination
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(skip).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find expenses: %w", err)
	}
	defer cursor.Close(ctx)

	var expenses []*domain.Expense
	if err := cursor.All(ctx, &expenses); err != nil {
		return nil, 0, fmt.Errorf("failed to decode expenses: %w", err)
	}

	return expenses, total, nil
}

func (r *expenseRepository) Update(ctx context.Context, expense *domain.Expense) error {
	expense.UpdatedAt = time.Now()

	update := bson.M{
		"$set": expense,
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": expense.ID}, update)
	if err != nil {
		return fmt.Errorf("failed to update expense: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("expense not found")
	}

	return nil
}

func (r *expenseRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid expense ID: %w", err)
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete expense: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("expense not found")
	}

	return nil
}

func (r *expenseRepository) UpdateStatus(ctx context.Context, id string, status domain.ExpenseStatus) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid expense ID: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return fmt.Errorf("failed to update expense status: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("expense not found")
	}

	return nil
}

func (r *expenseRepository) FindPendingByCompanyID(ctx context.Context, companyID string) ([]*domain.Expense, error) {
	objectID, err := primitive.ObjectIDFromHex(companyID)
	if err != nil {
		return nil, fmt.Errorf("invalid company ID: %w", err)
	}

	filter := bson.M{
		"company_id": objectID,
		"status":     domain.StatusPending,
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find pending expenses: %w", err)
	}
	defer cursor.Close(ctx)

	var expenses []*domain.Expense
	if err := cursor.All(ctx, &expenses); err != nil {
		return nil, fmt.Errorf("failed to decode expenses: %w", err)
	}

	return expenses, nil
}
