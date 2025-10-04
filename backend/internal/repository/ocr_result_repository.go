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

type ocrResultRepository struct {
	collection *mongo.Collection
}

// NewOCRResultRepository creates a new OCR result repository
func NewOCRResultRepository() domain.OCRResultRepository {
	return &ocrResultRepository{
		collection: database.GetCollection("ocr_results"),
	}
}

func (r *ocrResultRepository) Create(ctx context.Context, result *domain.OCRResult) error {
	result.CreatedAt = time.Now()
	result.ProcessedAt = time.Now()

	insertResult, err := r.collection.InsertOne(ctx, result)
	if err != nil {
		return fmt.Errorf("failed to create OCR result: %w", err)
	}

	result.ID = insertResult.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *ocrResultRepository) FindByID(ctx context.Context, id string) (*domain.OCRResult, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid OCR result ID: %w", err)
	}

	var result domain.OCRResult
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("OCR result not found")
		}
		return nil, fmt.Errorf("failed to find OCR result: %w", err)
	}

	return &result, nil
}

func (r *ocrResultRepository) FindByReceiptURL(ctx context.Context, receiptURL string) (*domain.OCRResult, error) {
	var result domain.OCRResult
	err := r.collection.FindOne(ctx, bson.M{"receipt_url": receiptURL}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("OCR result not found")
		}
		return nil, fmt.Errorf("failed to find OCR result: %w", err)
	}

	return &result, nil
}

func (r *ocrResultRepository) FindByUserID(ctx context.Context, userID string) ([]*domain.OCRResult, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	cursor, err := r.collection.Find(ctx, bson.M{"user_id": objectID})
	if err != nil {
		return nil, fmt.Errorf("failed to find OCR results: %w", err)
	}
	defer cursor.Close(ctx)

	var results []*domain.OCRResult
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode OCR results: %w", err)
	}

	return results, nil
}
