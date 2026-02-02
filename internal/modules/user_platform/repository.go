package user_platform

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserPlatformRepository struct {
	userPlatforms *mongo.Collection
}

func NewUserPlatformRepository(db *mongo.Database) *UserPlatformRepository {
	return &UserPlatformRepository{
		userPlatforms: db.Collection("user_platforms"),
	}
}

func (r *UserPlatformRepository) CreateUserPlatform(ctx context.Context, userPlatform *UserPlatform) error {
	userPlatform.ID = primitive.NewObjectID()
	userPlatform.CreatedAt = time.Now()
	userPlatform.UpdatedAt = time.Now()
	_, err := r.userPlatforms.InsertOne(ctx, userPlatform)
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
