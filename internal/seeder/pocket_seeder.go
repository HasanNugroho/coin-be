package seeder

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PocketSeeder struct {
	db *mongo.Database
}

func NewPocketSeeder(db *mongo.Database) *PocketSeeder {
	return &PocketSeeder{db: db}
}

func (ps *PocketSeeder) SeedMainPockets(ctx context.Context) error {
	usersCollection := ps.db.Collection("users")
	pocketsCollection := ps.db.Collection("pockets")

	// Get all users
	cursor, err := usersCollection.Find(ctx, bson.M{"deleted_at": nil})
	if err != nil {
		return fmt.Errorf("failed to fetch users: %w", err)
	}
	defer cursor.Close(ctx)

	var users []bson.M
	if err = cursor.All(ctx, &users); err != nil {
		return fmt.Errorf("failed to decode users: %w", err)
	}

	now := time.Now()
	createdCount := 0

	for _, user := range users {
		userID := user["_id"].(primitive.ObjectID)

		// Check if user already has a MAIN pocket
		existingMainPocket := pocketsCollection.FindOne(ctx, bson.M{
			"user_id": userID,
			"type":    "main",
			"deleted_at": nil,
		})

		if existingMainPocket.Err() == nil {
			// Main pocket already exists, skip
			continue
		}

		if existingMainPocket.Err() != mongo.ErrNoDocuments {
			return fmt.Errorf("error checking main pocket for user %s: %w", userID.Hex(), existingMainPocket.Err())
		}

		// Create MAIN pocket for user
		mainPocket := bson.M{
			"_id":              primitive.NewObjectID(),
			"user_id":          userID,
			"name":             "Main Pocket",
			"type":             "main",
			"category_id":      nil,
			"balance":          0,
			"is_default":       true,
			"is_active":        true,
			"is_locked":        false,
			"icon":             nil,
			"icon_color":       nil,
			"background_color": nil,
			"created_at":       now,
			"updated_at":       now,
			"deleted_at":       nil,
		}

		_, err := pocketsCollection.InsertOne(ctx, mainPocket)
		if err != nil {
			return fmt.Errorf("failed to create main pocket for user %s: %w", userID.Hex(), err)
		}

		createdCount++
	}

	fmt.Printf("Successfully created %d MAIN pockets for users\n", createdCount)
	return nil
}
