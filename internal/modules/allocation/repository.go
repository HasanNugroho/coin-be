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
	allocations    *mongo.Collection
	allocationLogs *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		allocations:    db.Collection("allocations"),
		allocationLogs: db.Collection("allocation_logs"),
	}
}

func (r *Repository) Create(ctx context.Context, allocation *Allocation) error {
	allocation.ID = primitive.NewObjectID()
	allocation.CreatedAt = time.Now()
	_, err := r.allocations.InsertOne(ctx, allocation)
	return err
}

func (r *Repository) GetByID(ctx context.Context, id primitive.ObjectID) (*Allocation, error) {
	var allocation Allocation
	err := r.allocations.FindOne(ctx, bson.M{"_id": id}).Decode(&allocation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("allocation not found")
		}
		return nil, err
	}
	return &allocation, nil
}

func (r *Repository) GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]*Allocation, error) {
	cursor, err := r.allocations.Find(ctx, bson.M{"user_id": userID})
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

func (r *Repository) GetActiveByUserID(ctx context.Context, userID primitive.ObjectID) ([]*Allocation, error) {
	opts := options.Find().SetSort(bson.D{{Key: "priority", Value: 1}})
	cursor, err := r.allocations.Find(ctx, bson.M{"user_id": userID, "is_active": true}, opts)
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

func (r *Repository) Update(ctx context.Context, id primitive.ObjectID, allocation *Allocation) error {
	result, err := r.allocations.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": allocation})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("allocation not found")
	}
	return nil
}

func (r *Repository) UpdateCurrentAmount(ctx context.Context, id primitive.ObjectID, amount float64) error {
	result, err := r.allocations.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$inc": bson.M{"current_amount": amount}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("allocation not found")
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.allocations.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("allocation not found")
	}
	return nil
}

func (r *Repository) CreateLog(ctx context.Context, log *AllocationLog) error {
	log.ID = primitive.NewObjectID()
	log.CreatedAt = time.Now()
	_, err := r.allocationLogs.InsertOne(ctx, log)
	return err
}

func (r *Repository) GetLogsByUserID(ctx context.Context, userID primitive.ObjectID, limit, skip int64) ([]*AllocationLog, error) {
	opts := options.Find().
		SetLimit(limit).
		SetSkip(skip).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.allocationLogs.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*AllocationLog
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}

func (r *Repository) GetLogsByAllocationID(ctx context.Context, allocationID primitive.ObjectID, limit, skip int64) ([]*AllocationLog, error) {
	opts := options.Find().
		SetLimit(limit).
		SetSkip(skip).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.allocationLogs.Find(ctx, bson.M{"allocation_id": allocationID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*AllocationLog
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}

func (r *Repository) GetLogsByTransactionID(ctx context.Context, transactionID primitive.ObjectID) ([]*AllocationLog, error) {
	cursor, err := r.allocationLogs.Find(ctx, bson.M{"transaction_id": transactionID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*AllocationLog
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}

func (r *Repository) CreateDefaultAllocations(ctx context.Context, userID primitive.ObjectID) error {
	defaultAllocations := []Allocation{
		{UserID: userID, Name: "Bills & Utilities", Priority: 1, Percentage: 40, CurrentAmount: 0, IsActive: true},
		{UserID: userID, Name: "Emergency Fund", Priority: 2, Percentage: 10, CurrentAmount: 0, TargetAmount: floatPtr(10000000), IsActive: true},
		{UserID: userID, Name: "Investment", Priority: 3, Percentage: 30, CurrentAmount: 0, IsActive: true},
		{UserID: userID, Name: "Savings", Priority: 4, Percentage: 20, CurrentAmount: 0, TargetAmount: floatPtr(5000000), IsActive: true},
	}

	for _, alloc := range defaultAllocations {
		if err := r.Create(ctx, &alloc); err != nil {
			return err
		}
	}
	return nil
}

func floatPtr(f float64) *float64 {
	return &f
}

func (r *Repository) GetTotalCurrentAmount(ctx context.Context, userID primitive.ObjectID) (float64, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"user_id": userID, "is_active": true}}},
		{{Key: "$group", Value: bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$current_amount"},
		}}},
	}

	cursor, err := r.allocations.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		return 0, err
	}

	if len(result) == 0 {
		return 0, nil
	}

	total, ok := result[0]["total"].(float64)
	if !ok {
		return 0, nil
	}

	return total, nil
}
