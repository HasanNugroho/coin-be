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

// GetAllocationsByExecuteDayWithOverflow fetches allocations for a specific day, including overflow allocations
// that should execute on the last day of the month if their execute_day exceeds the month's length
func (r *Repository) GetAllocationsByExecuteDayWithOverflow(ctx context.Context, executeDay int, lastDayOfMonth int) ([]map[string]interface{}, error) {
	// Build match conditions
	matchConditions := []bson.D{
		{{Key: "execute_day", Value: executeDay}},
	}

	// If today is the last day of the month, also include allocations scheduled for days beyond this month
	if executeDay == lastDayOfMonth && lastDayOfMonth < 31 {
		matchConditions = append(matchConditions, bson.D{
			{Key: "execute_day", Value: bson.D{{Key: "$gt", Value: lastDayOfMonth}}},
		})
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "$or", Value: matchConditions},
			{Key: "is_active", Value: true},
			{Key: "deleted_at", Value: nil},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "user_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "user"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$user"},
			{Key: "preserveNullAndEmptyArrays", Value: false},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "user_profiles"},
			{Key: "localField", Value: "user_id"},
			{Key: "foreignField", Value: "user_id"},
			{Key: "as", Value: "user_profile"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$user_profile"},
			{Key: "preserveNullAndEmptyArrays", Value: false},
		}}},
		{{Key: "$match", Value: bson.D{
			{Key: "user.is_active", Value: true},
			{Key: "user_profile.is_active", Value: true},
		}}},
		// Normalize execute_day for overflow allocations
		{{Key: "$addFields", Value: bson.D{
			{Key: "execute_day", Value: bson.D{
				{Key: "$cond", Value: bson.A{
					bson.D{{Key: "$gt", Value: bson.A{"$execute_day", lastDayOfMonth}}},
					lastDayOfMonth,
					"$execute_day",
				}},
			}},
		}}},
	}

	cursor, err := r.allocations.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []map[string]interface{}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// GetAllocationsByExecuteDay is kept for backward compatibility or simpler queries
func (r *Repository) GetAllocationsByExecuteDay(ctx context.Context, executeDay int, lastDayOfMonth int) ([]*Allocation, error) {
	matchConditions := bson.A{
		bson.M{"execute_day": executeDay},
	}

	// Include overflow allocations if today is the last day of a short month
	if executeDay == lastDayOfMonth && lastDayOfMonth < 31 {
		matchConditions = append(matchConditions, bson.M{
			"execute_day": bson.M{"$gt": lastDayOfMonth},
		})
	}

	cursor, err := r.allocations.Find(ctx, bson.M{
		"$or":        matchConditions,
		"is_active":  true,
		"deleted_at": nil,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var allocations []*Allocation
	if err = cursor.All(ctx, &allocations); err != nil {
		return nil, err
	}

	// Normalize execute_day for overflow allocations
	for _, alloc := range allocations {
		if alloc.ExecuteDay != nil && *alloc.ExecuteDay > lastDayOfMonth {
			alloc.ExecuteDay = &lastDayOfMonth
		}
	}

	return allocations, nil
}
