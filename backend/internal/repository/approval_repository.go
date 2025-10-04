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
	fmt.Printf("💾 Creating approval record - Expense: %s, Approver: %s, Status: %s\n",
		approval.ExpenseID.Hex(), approval.ApproverID.Hex(), approval.Status)

	approval.CreatedAt = time.Now()
	approval.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, approval)
	if err != nil {
		fmt.Printf("❌ Failed to insert approval: %v\n", err)
		return fmt.Errorf("failed to create approval: %w", err)
	}

	approval.ID = result.InsertedID.(primitive.ObjectID)
	fmt.Printf("✅ Approval created with ID: %s\n", approval.ID.Hex())

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
	fmt.Printf("🔍 FindPendingByApproverID - Looking for approver ID: %s\n", approverID)

	objectID, err := primitive.ObjectIDFromHex(approverID)
	if err != nil {
		fmt.Printf("❌ Invalid approver ID format: %v\n", err)
		return nil, fmt.Errorf("invalid approver ID: %w", err)
	}

	filter := bson.M{
		"approver_id": objectID,
		"status":      domain.ApprovalPending,
	}

	fmt.Printf("🔍 Query filter: %+v\n", filter)

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		fmt.Printf("❌ Database query failed: %v\n", err)
		return nil, fmt.Errorf("failed to find pending approvals: %w", err)
	}
	defer cursor.Close(ctx)

	var approvals []*domain.Approval
	if err := cursor.All(ctx, &approvals); err != nil {
		fmt.Printf("❌ Failed to decode approvals: %v\n", err)
		return nil, fmt.Errorf("failed to decode approvals: %w", err)
	}

	fmt.Printf("✅ Found %d approval records\n", len(approvals))
	for i, approval := range approvals {
		fmt.Printf("   [%d] Approval ID: %s, Expense ID: %s, Approver ID: %s, Status: %s\n",
			i+1, approval.ID.Hex(), approval.ExpenseID.Hex(), approval.ApproverID.Hex(), approval.Status)
	}

	return approvals, nil
}

// FindPendingByApproverIDWithDetails returns pending approvals with populated expense and user data
func (r *approvalRepository) FindPendingByApproverIDWithDetails(ctx context.Context, approverID string) ([]*domain.ApprovalWithDetails, error) {
	fmt.Printf("🔍 FindPendingByApproverIDWithDetails - Looking for approver ID: %s\n", approverID)

	objectID, err := primitive.ObjectIDFromHex(approverID)
	if err != nil {
		fmt.Printf("❌ Invalid approver ID format: %v\n", err)
		return nil, fmt.Errorf("invalid approver ID: %w", err)
	}

	// Build aggregation pipeline to join with expenses and users
	pipeline := []bson.M{
		// Match pending approvals for this approver
		{
			"$match": bson.M{
				"approver_id": objectID,
				"status":      domain.ApprovalPending,
			},
		},
		// Lookup expense details
		{
			"$lookup": bson.M{
				"from":         "expenses",
				"localField":   "expense_id",
				"foreignField": "_id",
				"as":           "expense_data",
			},
		},
		// Unwind expense array (should be single document)
		{
			"$unwind": bson.M{
				"path":                       "$expense_data",
				"preserveNullAndEmptyArrays": true,
			},
		},
		// Lookup user details from expense
		{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "expense_data.user_id",
				"foreignField": "_id",
				"as":           "expense_user",
			},
		},
		// Unwind user array
		{
			"$unwind": bson.M{
				"path":                       "$expense_user",
				"preserveNullAndEmptyArrays": true,
			},
		},
		// Lookup approver details
		{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "approver_id",
				"foreignField": "_id",
				"as":           "approver_data",
			},
		},
		// Unwind approver array
		{
			"$unwind": bson.M{
				"path":                       "$approver_data",
				"preserveNullAndEmptyArrays": true,
			},
		},
		// Project final structure
		{
			"$project": bson.M{
				"_id":         1,
				"expense_id":  1,
				"approver_id": 1,
				"level":       1,
				"status":      1,
				"comments":    1,
				"approved_at": 1,
				"created_at":  1,
				"updated_at":  1,
				"expense": bson.M{
					"_id":                    "$expense_data._id",
					"user_id":                "$expense_data.user_id",
					"company_id":             "$expense_data.company_id",
					"amount":                 "$expense_data.amount",
					"currency":               "$expense_data.currency",
					"converted_amount":       "$expense_data.converted_amount",
					"exchange_rate":          "$expense_data.exchange_rate",
					"category":               "$expense_data.category",
					"description":            "$expense_data.description",
					"expense_date":           "$expense_data.expense_date",
					"receipt_url":            "$expense_data.receipt_url",
					"merchant":               "$expense_data.merchant",
					"status":                 "$expense_data.status",
					"current_approval_level": "$expense_data.current_approval_level",
					"created_at":             "$expense_data.created_at",
					"updated_at":             "$expense_data.updated_at",
					"user":                   "$expense_user",
				},
				"approver": "$approver_data",
			},
		},
	}

	fmt.Printf("🔍 Executing aggregation pipeline with %d stages\n", len(pipeline))

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		fmt.Printf("❌ Aggregation query failed: %v\n", err)
		return nil, fmt.Errorf("failed to aggregate pending approvals: %w", err)
	}
	defer cursor.Close(ctx)

	var approvals []*domain.ApprovalWithDetails
	if err := cursor.All(ctx, &approvals); err != nil {
		fmt.Printf("❌ Failed to decode approvals with details: %v\n", err)
		return nil, fmt.Errorf("failed to decode approvals: %w", err)
	}

	fmt.Printf("✅ Found %d approval records with details\n", len(approvals))
	for i, approval := range approvals {
		expenseDesc := "nil"
		userName := "nil"
		if approval.Expense != nil {
			expenseDesc = approval.Expense.Description
			if approval.Expense.User != nil {
				userName = approval.Expense.User.FirstName + " " + approval.Expense.User.LastName
			}
		}
		fmt.Printf("   [%d] Approval ID: %s, Expense: %s, User: %s, Status: %s\n",
			i+1, approval.ID.Hex(), expenseDesc, userName, approval.Status)
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
