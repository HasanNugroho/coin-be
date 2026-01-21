package target

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
	targets *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		targets: db.Collection("saving_targets"),
	}
}

func (r *Repository) Create(ctx context.Context, target *SavingTarget) error {
	target.ID = primitive.NewObjectID()
	target.CreatedAt = time.Now()
	_, err := r.targets.InsertOne(ctx, target)
	return err
}

func (r *Repository) GetByID(ctx context.Context, id primitive.ObjectID) (*SavingTarget, error) {
	var target SavingTarget
	err := r.targets.FindOne(ctx, bson.M{"_id": id}).Decode(&target)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("saving target not found")
		}
		return nil, err
	}
	return &target, nil
}

func (r *Repository) GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]*SavingTarget, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.targets.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var targets []*SavingTarget
	if err = cursor.All(ctx, &targets); err != nil {
		return nil, err
	}
	return targets, nil
}

func (r *Repository) GetByAllocationID(ctx context.Context, allocationID primitive.ObjectID) ([]*SavingTarget, error) {
	cursor, err := r.targets.Find(ctx, bson.M{"allocation_id": allocationID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var targets []*SavingTarget
	if err = cursor.All(ctx, &targets); err != nil {
		return nil, err
	}
	return targets, nil
}

func (r *Repository) GetActiveByUserID(ctx context.Context, userID primitive.ObjectID) ([]*SavingTarget, error) {
	cursor, err := r.targets.Find(ctx, bson.M{"user_id": userID, "status": TargetStatusActive})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var targets []*SavingTarget
	if err = cursor.All(ctx, &targets); err != nil {
		return nil, err
	}
	return targets, nil
}

func (r *Repository) Update(ctx context.Context, id primitive.ObjectID, target *SavingTarget) error {
	result, err := r.targets.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": target})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("saving target not found")
	}
	return nil
}

func (r *Repository) UpdateCurrentAmount(ctx context.Context, id primitive.ObjectID, amount float64) error {
	result, err := r.targets.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"current_amount": amount}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("saving target not found")
	}
	return nil
}

func (r *Repository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status TargetStatus) error {
	result, err := r.targets.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"status": status}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("saving target not found")
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.targets.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("saving target not found")
	}
	return nil
}
