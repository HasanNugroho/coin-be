package platform

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	platforms *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		platforms: db.Collection("platforms"),
	}
}

func (r *Repository) CreatePlatform(ctx context.Context, platform *Platform) error {
	platform.ID = primitive.NewObjectID()
	platform.CreatedAt = time.Now()
	platform.UpdatedAt = time.Now()
	_, err := r.platforms.InsertOne(ctx, platform)
	return err
}

func (r *Repository) GetPlatformByID(ctx context.Context, id primitive.ObjectID) (*Platform, error) {
	var platform Platform
	err := r.platforms.FindOne(ctx, bson.M{"_id": id, "deleted_at": nil}).Decode(&platform)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("platform not found")
		}
		return nil, err
	}
	return &platform, nil
}

func (r *Repository) GetPlatformByName(ctx context.Context, name string) (*Platform, error) {
	var platform Platform
	err := r.platforms.FindOne(ctx, bson.M{"name": name, "deleted_at": nil}).Decode(&platform)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("platform not found")
		}
		return nil, err
	}
	return &platform, nil
}

func (r *Repository) UpdatePlatform(ctx context.Context, id primitive.ObjectID, platform *Platform) error {
	platform.UpdatedAt = time.Now()
	result, err := r.platforms.UpdateOne(ctx, bson.M{"_id": id, "deleted_at": nil}, bson.M{"$set": platform})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("platform not found")
	}
	return nil
}

func (r *Repository) DeletePlatform(ctx context.Context, id primitive.ObjectID) error {
	now := time.Now()
	result, err := r.platforms.UpdateOne(
		ctx,
		bson.M{"_id": id, "deleted_at": nil},
		bson.M{
			"$set": bson.M{
				"deleted_at": now,
				"updated_at": now,
			},
		},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("platform not found")
	}
	return nil
}

func (r *Repository) ListPlatforms(ctx context.Context, limit int64, skip int64) ([]*Platform, error) {
	opts := options.Find().SetLimit(limit).SetSkip(skip).SetSort(bson.M{"created_at": -1})
	cursor, err := r.platforms.Find(ctx, bson.M{"deleted_at": nil}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var platforms []*Platform
	if err = cursor.All(ctx, &platforms); err != nil {
		return nil, err
	}
	return platforms, nil
}

func (r *Repository) ListActivePlatforms(ctx context.Context, limit int64, skip int64) ([]*Platform, error) {
	opts := options.Find().SetLimit(limit).SetSkip(skip).SetSort(bson.M{"created_at": -1})
	cursor, err := r.platforms.Find(ctx, bson.M{"is_active": true, "deleted_at": nil}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var platforms []*Platform
	if err = cursor.All(ctx, &platforms); err != nil {
		return nil, err
	}
	return platforms, nil
}

func (r *Repository) ListPlatformsByType(ctx context.Context, platformType string, limit int64, skip int64) ([]*Platform, error) {
	opts := options.Find().SetLimit(limit).SetSkip(skip).SetSort(bson.M{"created_at": -1})
	cursor, err := r.platforms.Find(ctx, bson.M{"type": platformType, "deleted_at": nil}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var platforms []*Platform
	if err = cursor.All(ctx, &platforms); err != nil {
		return nil, err
	}
	return platforms, nil
}

func (r *Repository) CountPlatforms(ctx context.Context) (int64, error) {
	return r.platforms.CountDocuments(ctx, bson.M{"deleted_at": nil})
}
