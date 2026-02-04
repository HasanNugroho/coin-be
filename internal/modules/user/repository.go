package user

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
	users        *mongo.Collection
	userProfiles *mongo.Collection
	roles        *mongo.Collection
	userRoles    *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		users:        db.Collection("users"),
		userProfiles: db.Collection("user_profiles"),
		roles:        db.Collection("roles"),
		userRoles:    db.Collection("user_roles"),
	}
}

func (r *Repository) CreateUser(ctx context.Context, user *User) error {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	_, err := r.users.InsertOne(ctx, user)
	return err
}

func (r *Repository) GetUserByID(ctx context.Context, id primitive.ObjectID) (*User, error) {
	var user User
	err := r.users.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.users.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) UpdateUser(ctx context.Context, id primitive.ObjectID, user *User) error {
	user.UpdatedAt = time.Now()
	result, err := r.users.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": user})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *Repository) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.users.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *Repository) ListUsers(ctx context.Context, limit int64, skip int64, role, search, sort, order string) ([]*User, int64, error) {
	filter := bson.M{}

	if role != "" {
		filter["role"] = role
	}

	if search != "" {
		filter["$or"] = bson.A{
			bson.M{"name": bson.M{"$regex": search, "$options": "i"}},
			bson.M{"email": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	sortField := "created_at"
	if sort != "" {
		switch sort {
		case "name", "email", "created_at", "updated_at":
			sortField = sort
		}
	}

	sortOrder := int32(-1)
	if order == "asc" {
		sortOrder = 1
	}

	opts := options.Find().SetLimit(limit).SetSkip(skip).SetSort(bson.M{sortField: sortOrder})
	cursor, err := r.users.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var users []*User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}

	total, err := r.users.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *Repository) CreateUserProfile(ctx context.Context, profile *UserProfile) error {
	profile.ID = primitive.NewObjectID()
	profile.CreatedAt = time.Now()
	profile.UpdatedAt = time.Now()
	_, err := r.userProfiles.InsertOne(ctx, profile)
	return err
}

func (r *Repository) GetUserProfileByUserID(ctx context.Context, userID primitive.ObjectID) (*UserProfile, error) {
	var profile UserProfile
	err := r.userProfiles.FindOne(ctx, bson.M{"user_id": userID}).Decode(&profile)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user profile not found")
		}
		return nil, err
	}
	return &profile, nil
}

func (r *Repository) UpdateUserProfile(ctx context.Context, userID primitive.ObjectID, profile *UserProfile) error {
	profile.UpdatedAt = time.Now()
	result, err := r.userProfiles.UpdateOne(ctx, bson.M{"user_id": userID}, bson.M{"$set": profile})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("user profile not found")
	}
	return nil
}

func (r *Repository) DeleteUserProfile(ctx context.Context, userID primitive.ObjectID) error {
	result, err := r.userProfiles.DeleteOne(ctx, bson.M{"user_id": userID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("user profile not found")
	}
	return nil
}

func (r *Repository) CountUsersByRole(ctx context.Context, role string) (int64, error) {
	count, err := r.users.CountDocuments(ctx, bson.M{"role": role})
	return count, err
}
