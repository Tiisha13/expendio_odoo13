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

type companyRepository struct {
	collection *mongo.Collection
}

// NewCompanyRepository creates a new company repository
func NewCompanyRepository() domain.CompanyRepository {
	return &companyRepository{
		collection: database.GetCollection("companies"),
	}
}

func (r *companyRepository) Create(ctx context.Context, company *domain.Company) error {
	company.CreatedAt = time.Now()
	company.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, company)
	if err != nil {
		return fmt.Errorf("failed to create company: %w", err)
	}

	company.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *companyRepository) FindByID(ctx context.Context, id string) (*domain.Company, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid company ID: %w", err)
	}

	var company domain.Company
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&company)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("company not found")
		}
		return nil, fmt.Errorf("failed to find company: %w", err)
	}

	return &company, nil
}

func (r *companyRepository) FindByName(ctx context.Context, name string) (*domain.Company, error) {
	var company domain.Company
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&company)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("company not found")
		}
		return nil, fmt.Errorf("failed to find company: %w", err)
	}

	return &company, nil
}

func (r *companyRepository) Update(ctx context.Context, company *domain.Company) error {
	company.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"name":             company.Name,
			"base_currency":    company.BaseCurrency,
			"country":          company.Country,
			"admin_user_id":    company.AdminUserID,
			"approval_rule_id": company.ApprovalRuleID,
			"is_active":        company.IsActive,
			"updated_at":       company.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": company.ID}, update)
	if err != nil {
		return fmt.Errorf("failed to update company: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("company not found")
	}

	return nil
}

func (r *companyRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid company ID: %w", err)
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete company: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("company not found")
	}

	return nil
}
