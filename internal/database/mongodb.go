// internal/database/mongodb.go
package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"gorbit/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongoDB(cfg *config.Config) (*mongo.Client, error) {
	// Create MongoDB connection string
	connectionString := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=%s&authMechanism=%s",
		cfg.Databases.MongoDB.Username,
		cfg.Databases.MongoDB.Password,
		cfg.Databases.MongoDB.Host,
		cfg.Databases.MongoDB.Port,
		cfg.Databases.MongoDB.Database,
		cfg.Databases.MongoDB.AuthSource,
		cfg.Databases.MongoDB.AuthMechanism,
	)

	// Set client options
	clientOptions := options.Client().ApplyURI(connectionString)

	// Set context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Printf("Failed to connect to MongoDB: %v", err)
		return nil, err
	}

	// Verify the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Failed to ping MongoDB: %v", err)
		return nil, err
	}

	return client, nil
}
