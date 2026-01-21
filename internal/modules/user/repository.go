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
	users              *mongo.Collection
	financialProfiles  *mongo.Collection
	roles              *mongo.Collection
	userRoles          *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		users:              db.Collection("users"),
		financialProfiles:  db.Collection("financial_profiles"),
		roles:              db.Collection("roles"),
		userRoles:          db.Collection("user_roles"),
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

func (r *Repository) ListUsers(ctx context.Context, limit int64, skip int64) ([]*User, error) {
	opts := options.Find().SetLimit(limit).SetSkip(skip)
	cursor, err := r.users.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) CreateFinancialProfile(ctx context.Context, profile *FinancialProfile) error {
	profile.ID = primitive.NewObjectID()
	profile.CreatedAt = time.Now()
	profile.UpdatedAt = time.Now()
	_, err := r.financialProfiles.InsertOne(ctx, profile)
	return err
}

func (r *Repository) GetFinancialProfileByUserID(ctx context.Context, userID primitive.ObjectID) (*FinancialProfile, error) {
	var profile FinancialProfile
	err := r.financialProfiles.FindOne(ctx, bson.M{"user_id": userID}).Decode(&profile)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("financial profile not found")
		}
		return nil, err
	}
	return &profile, nil
}

func (r *Repository) UpdateFinancialProfile(ctx context.Context, userID primitive.ObjectID, profile *FinancialProfile) error {
	profile.UpdatedAt = time.Now()
	result, err := r.financialProfiles.UpdateOne(ctx, bson.M{"user_id": userID}, bson.M{"$set": profile})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("financial profile not found")
	}
	return nil
}

func (r *Repository) DeleteFinancialProfile(ctx context.Context, userID primitive.ObjectID) error {
	result, err := r.financialProfiles.DeleteOne(ctx, bson.M{"user_id": userID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("financial profile not found")
	}
	return nil
}

func (r *Repository) CreateRole(ctx context.Context, role *Role) error {
	role.ID = primitive.NewObjectID()
	role.CreatedAt = time.Now()
	_, err := r.roles.InsertOne(ctx, role)
	return err
}

func (r *Repository) GetRoleByID(ctx context.Context, id primitive.ObjectID) (*Role, error) {
	var role Role
	err := r.roles.FindOne(ctx, bson.M{"_id": id}).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("role not found")
		}
		return nil, err
	}
	return &role, nil
}

func (r *Repository) GetRoleByName(ctx context.Context, name string) (*Role, error) {
	var role Role
	err := r.roles.FindOne(ctx, bson.M{"name": name}).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("role not found")
		}
		return nil, err
	}
	return &role, nil
}

func (r *Repository) ListRoles(ctx context.Context) ([]*Role, error) {
	cursor, err := r.roles.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var roles []*Role
	if err = cursor.All(ctx, &roles); err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *Repository) AssignRoleToUser(ctx context.Context, userRole *UserRole) error {
	userRole.ID = primitive.NewObjectID()
	userRole.AssignedAt = time.Now()
	_, err := r.userRoles.InsertOne(ctx, userRole)
	return err
}

func (r *Repository) GetUserRoles(ctx context.Context, userID primitive.ObjectID) ([]*UserRole, error) {
	cursor, err := r.userRoles.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var userRoles []*UserRole
	if err = cursor.All(ctx, &userRoles); err != nil {
		return nil, err
	}
	return userRoles, nil
}

func (r *Repository) RemoveRoleFromUser(ctx context.Context, userID, roleID primitive.ObjectID) error {
	result, err := r.userRoles.DeleteOne(ctx, bson.M{"user_id": userID, "role_id": roleID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("user role not found")
	}
	return nil
}
