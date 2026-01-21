package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/HasanNugroho/coin-be/internal/core/config"
	"github.com/HasanNugroho/coin-be/internal/seeder"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Verify connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB successfully")

	// Get database
	db := client.Database(cfg.MongoDB)

	// Run seeder
	s := seeder.NewSeeder(db)
	if err := s.Seed(ctx); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}

	fmt.Println("\n✅ Database seeding completed successfully!")
	fmt.Println("\nDefault data has been set up:")
	fmt.Println("- 12 default categories (4 income, 8 expense)")
	fmt.Println("- 4 default allocations (Bills, Emergency Fund, Investment, Savings)")
	fmt.Println("\nYou can now start using the application!")
}
