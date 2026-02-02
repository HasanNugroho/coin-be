package allocation

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
	allocations *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		allocations: db.Collection("allocations"),
	}
}

func (r *Repository) CreateAllocation(ctx context.Context, allocation *Allocation) error {
	allocation.ID = primitive.NewObjectID()
	allocation.CreatedAt = time.Now()
	allocation.UpdatedAt = time.Now()
	_, err := r.allocations.InsertOne(ctx, allocation)
	return err
}

func (r *Repository) GetAllocationByID(ctx context.Context, id primitive.ObjectID) (*Allocation, error) {
	var allocation Allocation
	err := r.allocations.FindOne(ctx, bson.M{"_id": id, "deleted_at": nil}).Decode(&allocation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("allocation not found")
		}
		return nil, err
	}
	return &allocation, nil
}

func (r *Repository) GetAllocationsByUserID(ctx context.Context, userID primitive.ObjectID) ([]*Allocation, error) {
	cursor, err := r.allocations.Find(ctx, bson.M{"user_id": userID, "deleted_at": nil})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var allocations []*Allocation
	if err = cursor.All(ctx, &allocations); err != nil {
		return nil, err
	}
	return allocations, nil
}

func (r *Repository) GetActiveAllocationsByUserID(ctx context.Context, userID primitive.ObjectID) ([]*Allocation, error) {
	opts := options.Find().SetSort(bson.D{{Key: "priority", Value: 1}})
	cursor, err := r.allocations.Find(ctx, bson.M{
		"user_id":    userID,
		"is_active":  true,
		"deleted_at": nil,
	}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var allocations []*Allocation
	if err = cursor.All(ctx, &allocations); err != nil {
		return nil, err
	}
	return allocations, nil
}

func (r *Repository) UpdateAllocation(ctx context.Context, id primitive.ObjectID, allocation *Allocation) error {
	allocation.UpdatedAt = time.Now()
	result, err := r.allocations.UpdateOne(
		ctx,
		bson.M{"_id": id, "deleted_at": nil},
		bson.M{"$set": allocation},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("allocation not found")
	}
	return nil
}

func (r *Repository) DeleteAllocation(ctx context.Context, id primitive.ObjectID) error {
	now := time.Now()
	result, err := r.allocations.UpdateOne(
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
		return errors.New("allocation not found")
	}
	return nil
}
