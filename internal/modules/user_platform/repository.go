package user_platform

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserPlatformRepository struct {
	userPlatforms *mongo.Collection
}

func NewUserPlatformRepository(db *mongo.Database) *UserPlatformRepository {
	return &UserPlatformRepository{
		userPlatforms: db.Collection("user_platforms"),
	}
}

func (r *UserPlatformRepository) EnsureIndexes(ctx context.Context) {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "deleted_at", Value: 1},
				{Key: "is_active", Value: 1},
				{Key: "last_use_at", Value: -1},
			},
			Options: options.Index().SetName("idx_user_platforms_dropdown"),
		},
	}

	_, err := r.userPlatforms.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		log.Printf("failed to create user platform indexes: %v", err)
	}
}

func (r *UserPlatformRepository) CreateUserPlatform(ctx context.Context, userPlatform *UserPlatform) error {
	userPlatform.ID = primitive.NewObjectID()
	userPlatform.CreatedAt = time.Now()
	userPlatform.UpdatedAt = time.Now()
	_, err := r.userPlatforms.InsertOne(ctx, userPlatform)
	return err
}

func (r *UserPlatformRepository) CreateUserPlatformBulk(ctx context.Context, userPlatforms []*UserPlatform) error {
	for _, userPlatform := range userPlatforms {
		userPlatform.ID = primitive.NewObjectID()
		userPlatform.CreatedAt = time.Now()
		userPlatform.UpdatedAt = time.Now()
	}

	// Convert []*UserPlatform to []interface{}
	docs := make([]interface{}, len(userPlatforms))
	for i, userPlatform := range userPlatforms {
		docs[i] = userPlatform
	}

	_, err := r.userPlatforms.InsertMany(ctx, docs)
	return err
}

func (r *UserPlatformRepository) GetUserPlatformByID(ctx context.Context, id primitive.ObjectID) (*UserPlatform, error) {
	var userPlatform UserPlatform
	err := r.userPlatforms.FindOne(ctx, bson.M{"_id": id, "deleted_at": nil}).Decode(&userPlatform)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user platform not found")
		}
		return nil, err
	}
	return &userPlatform, nil
}

func (r *UserPlatformRepository) GetUserPlatformByUserAndPlatform(ctx context.Context, userID, platformID primitive.ObjectID) (*UserPlatform, error) {
	var userPlatform UserPlatform
	err := r.userPlatforms.FindOne(ctx, bson.M{
		"user_id":     userID,
		"platform_id": platformID,
		"deleted_at":  nil,
	}).Decode(&userPlatform)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user platform not found")
		}
		return nil, err
	}
	return &userPlatform, nil
}

func (r *UserPlatformRepository) GetUserPlatformsByUserID(ctx context.Context, userID primitive.ObjectID) ([]*UserPlatform, error) {
	cursor, err := r.userPlatforms.Find(ctx, bson.M{"user_id": userID, "deleted_at": nil})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var userPlatforms []*UserPlatform
	if err = cursor.All(ctx, &userPlatforms); err != nil {
		return nil, err
	}
	return userPlatforms, nil
}

func (r *UserPlatformRepository) GetUserPlatformsByUserIDDropdown(ctx context.Context, userID primitive.ObjectID) ([]*UserPlatform, error) {
	opts := options.Find().SetSort(bson.D{{Key: "last_use_at", Value: -1}})

	cursor, err := r.userPlatforms.Find(ctx, bson.M{"user_id": userID, "deleted_at": nil, "is_active": true}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var userPlatforms []*UserPlatform
	if err = cursor.All(ctx, &userPlatforms); err != nil {
		return nil, err
	}
	return userPlatforms, nil
}

func (r *UserPlatformRepository) UpdateUserPlatform(ctx context.Context, id primitive.ObjectID, userPlatform *UserPlatform) error {
	userPlatform.UpdatedAt = time.Now()
	result, err := r.userPlatforms.UpdateOne(
		ctx,
		bson.M{"_id": id, "deleted_at": nil},
		bson.M{"$set": userPlatform},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("user platform not found")
	}
	return nil
}

func (r *UserPlatformRepository) DeleteUserPlatform(ctx context.Context, id primitive.ObjectID) error {
	now := time.Now()
	result, err := r.userPlatforms.UpdateOne(
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
		return errors.New("user platform not found")
	}
	return nil
}
