package seeder

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

type Seeder struct {
	db *mongo.Database
}

func NewSeeder(db *mongo.Database) *Seeder {
	return &Seeder{
		db: db,
	}
}

func (s *Seeder) Seed(ctx context.Context) error {
	log.Println("Starting database seeding...")

	if err := s.seedCategories(ctx); err != nil {
		return fmt.Errorf("error seeding categories: %w", err)
	}

	if err := s.seedAllocations(ctx); err != nil {
		return fmt.Errorf("error seeding allocations: %w", err)
	}

	log.Println("Database seeding completed successfully!")
	return nil
}

func (s *Seeder) seedCategories(ctx context.Context) error {
	log.Println("Seeding categories...")

	collection := s.db.Collection("categories")

	// Check if categories already exist
	count, err := collection.EstimatedDocumentCount(ctx)
	if err != nil {
		return err
	}

	if count > 0 {
		log.Println("Categories already exist, skipping...")
		return nil
	}

	categories := getDefaultCategories()
	documents := make([]interface{}, len(categories))
	for i, cat := range categories {
		documents[i] = cat
	}

	result, err := collection.InsertMany(ctx, documents)
	if err != nil {
		return err
	}

	log.Printf("Inserted %d categories\n", len(result.InsertedIDs))
	return nil
}

func (s *Seeder) seedAllocations(ctx context.Context) error {
	log.Println("Seeding allocations...")

	collection := s.db.Collection("allocations")

	// Check if allocations already exist
	count, err := collection.EstimatedDocumentCount(ctx)
	if err != nil {
		return err
	}

	if count > 0 {
		log.Println("Allocations already exist, skipping...")
		return nil
	}

	allocations := getDefaultAllocations()
	documents := make([]interface{}, len(allocations))
	for i, alloc := range allocations {
		documents[i] = alloc
	}

	result, err := collection.InsertMany(ctx, documents)
	if err != nil {
		return err
	}

	log.Printf("Inserted %d allocations\n", len(result.InsertedIDs))
	return nil
}
