package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"expensio-backend/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	Client   *mongo.Client
	Database *mongo.Database
)

// ConnectMongoDB establishes connection to MongoDB
func ConnectMongoDB(cfg *config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set client options
	clientOptions := options.Client().
		ApplyURI(cfg.MongoDB.URI).
		SetMaxPoolSize(100).
		SetMinPoolSize(10).
		SetMaxConnIdleTime(30 * time.Second).
		SetServerSelectionTimeout(5 * time.Second)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	Client = client
	Database = client.Database(cfg.MongoDB.Database)

	// Create indexes
	if err := createIndexes(ctx); err != nil {
		log.Printf("⚠️  Warning: Failed to create some indexes: %v", err)
	}

	return nil
}

// DisconnectMongoDB closes MongoDB connection
func DisconnectMongoDB() {
	if Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := Client.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		} else {
			log.Println("✅ Disconnected from MongoDB")
		}
	}
}

// GetCollection returns a MongoDB collection
func GetCollection(collectionName string) *mongo.Collection {
	return Database.Collection(collectionName)
}

// createIndexes creates necessary indexes for collections
func createIndexes(ctx context.Context) error {
	// Users collection indexes
	usersCollection := GetCollection("users")
	_, err := usersCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    map[string]interface{}{"email": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: map[string]interface{}{"company_id": 1},
		},
		{
			Keys: map[string]interface{}{"role": 1},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create users indexes: %w", err)
	}

	// Companies collection indexes
	companiesCollection := GetCollection("companies")
	_, err = companiesCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    map[string]interface{}{"name": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("failed to create companies indexes: %w", err)
	}

	// Expenses collection indexes
	expensesCollection := GetCollection("expenses")
	_, err = expensesCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: map[string]interface{}{"user_id": 1},
		},
		{
			Keys: map[string]interface{}{"company_id": 1},
		},
		{
			Keys: map[string]interface{}{"status": 1},
		},
		{
			Keys: map[string]interface{}{"created_at": -1},
		},
		{
			Keys: map[string]interface{}{"company_id": 1, "status": 1, "created_at": -1},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create expenses indexes: %w", err)
	}

	// Approvals collection indexes
	approvalsCollection := GetCollection("approvals")
	_, err = approvalsCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: map[string]interface{}{"expense_id": 1},
		},
		{
			Keys: map[string]interface{}{"approver_id": 1, "status": 1},
		},
		{
			Keys: map[string]interface{}{"status": 1},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create approvals indexes: %w", err)
	}

	// Approval Rules collection indexes
	approvalRulesCollection := GetCollection("approval_rules")
	_, err = approvalRulesCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]interface{}{"company_id": 1},
	})
	if err != nil {
		return fmt.Errorf("failed to create approval_rules indexes: %w", err)
	}

	log.Println("✅ Database indexes created successfully")
	return nil
}

// HealthCheck checks if MongoDB connection is alive
func HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return Client.Ping(ctx, readpref.Primary())
}
