package admin_dashboard

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	users        *mongo.Collection
	transactions *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		users:        db.Collection("users"),
		transactions: db.Collection("transactions"),
	}
}

func (r *Repository) GetTotalUsersCount(ctx context.Context) (int64, error) {
	return r.users.CountDocuments(ctx, bson.M{})
}

func (r *Repository) GetActiveUsersCount(ctx context.Context) (int64, error) {
	return r.users.CountDocuments(ctx, bson.M{"is_active": true})
}

func (r *Repository) GetTotalTransactionsCount(ctx context.Context, startDate, endDate time.Time) (int64, error) {
	filter := bson.M{
		"deleted_at": nil,
		"date": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}
	return r.transactions.CountDocuments(ctx, filter)
}

func (r *Repository) GetTotalTransactionVolume(ctx context.Context, startDate, endDate time.Time) (float64, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"deleted_at": nil,
			"date": bson.M{
				"$gte": startDate,
				"$lte": endDate,
			},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$amount"},
		}}},
	}

	cursor, err := r.transactions.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		Total float64 `bson:"total"`
	}
	if err = cursor.All(ctx, &results); err != nil {
		return 0, err
	}

	if len(results) == 0 {
		return 0, nil
	}
	return results[0].Total, nil
}

func (r *Repository) GetUserGrowth(ctx context.Context, startDate, endDate time.Time) ([]UserGrowthData, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"created_at": bson.M{
				"$gte": startDate,
				"$lte": endDate,
			},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id": bson.M{
				"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$created_at"},
			},
			"count": bson.M{"$sum": 1},
		}}},
		{{Key: "$sort", Value: bson.M{"_id": 1}}},
	}

	cursor, err := r.users.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		Date  string `bson:"_id"`
		Count int64  `bson:"count"`
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	growth := make([]UserGrowthData, len(results))
	for i, res := range results {
		growth[i] = UserGrowthData{
			Date:  res.Date,
			Count: res.Count,
		}
	}
	return growth, nil
}
