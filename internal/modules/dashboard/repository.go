package dashboard

import (
	"context"
	"time"

	"github.com/HasanNugroho/coin-be/internal/modules/daily_summary"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	transactions   *mongo.Collection
	pockets        *mongo.Collection
	userCategories *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		transactions:   db.Collection("transactions"),
		pockets:        db.Collection("pockets"),
		userCategories: db.Collection("user_categories"),
	}
}

func (r *Repository) GetLiveDeltaSummary(ctx context.Context, userID primitive.ObjectID, startDate time.Time) (float64, float64, []daily_summary.CategoryBreakdown, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"user_id": userID,
			"$or": bson.A{
				bson.M{"deleted_at": bson.M{"$exists": false}},
				bson.M{"deleted_at": nil},
			},
			"date": bson.M{"$gte": startDate.UTC()},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id": bson.M{
				"type":        "$type",
				"category_id": "$category_id",
			},
			"amount": bson.M{"$sum": "$amount"},
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "user_categories",
			"localField":   "_id.category_id",
			"foreignField": "_id",
			"as":           "category_info",
		}}},
		{{Key: "$project", Value: bson.M{
			"_id":    1,
			"amount": 1,
			"category_name": bson.M{
				"$ifNull": bson.A{
					bson.M{"$arrayElemAt": bson.A{"$category_info.name", 0}},
					"Uncategorized",
				},
			},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id": nil,
			"total_income": bson.M{
				"$sum": bson.M{
					"$cond": bson.A{
						bson.M{"$eq": bson.A{"$_id.type", "income"}},
						"$amount",
						0,
					},
				},
			},
			"total_expense": bson.M{
				"$sum": bson.M{
					"$cond": bson.A{
						bson.M{"$eq": bson.A{"$_id.type", "expense"}},
						"$amount",
						0,
					},
				},
			},
			"categories": bson.M{
				"$push": bson.M{
					"type":          "$_id.type",
					"category_id":   "$_id.category_id",
					"category_name": "$category_name",
					"amount":        "$amount",
				},
			},
		}}},
	}

	cursor, err := r.transactions.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, 0, nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		TotalIncome  float64 `bson:"total_income"`
		TotalExpense float64 `bson:"total_expense"`
		Categories   []struct {
			Type         string              `bson:"type"`
			CategoryID   *primitive.ObjectID `bson:"category_id"`
			CategoryName string              `bson:"category_name"`
			Amount       float64             `bson:"amount"`
		} `bson:"categories"`
	}

	if err = cursor.All(ctx, &results); err != nil {
		return 0, 0, nil, err
	}

	if len(results) == 0 {
		return 0, 0, []daily_summary.CategoryBreakdown{}, nil
	}

	categories := make([]daily_summary.CategoryBreakdown, 0, len(results[0].Categories))
	for _, cat := range results[0].Categories {
		categories = append(categories, daily_summary.CategoryBreakdown{
			CategoryID:   cat.CategoryID,
			CategoryName: cat.CategoryName,
			Type:         cat.Type,
			Amount:       cat.Amount,
		})
	}

	return results[0].TotalIncome, results[0].TotalExpense, categories, nil
}

func (r *Repository) GetTotalNetWorth(ctx context.Context, userID primitive.ObjectID) (float64, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"user_id":    userID,
			"is_active":  true,
			"deleted_at": nil,
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": bson.M{"$toDouble": "$balance"}},
		}}},
	}

	cursor, err := r.pockets.Aggregate(ctx, pipeline)
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

func (r *Repository) GetCategoryNames(ctx context.Context, ids []primitive.ObjectID) (map[primitive.ObjectID]string, error) {
	cursor, err := r.userCategories.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		ID   primitive.ObjectID `bson:"_id"`
		Name string             `bson:"name"`
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	nameMap := make(map[primitive.ObjectID]string)
	for _, res := range results {
		nameMap[res.ID] = res.Name
	}
	return nameMap, nil
}

func (r *Repository) EnsureIndexes(ctx context.Context) error {
	transactionIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "date", Value: -1},
			},
			Options: options.Index().
				SetName("idx_transactions_user_date"),
		},
	}

	_, err := r.transactions.Indexes().CreateMany(ctx, transactionIndexes)
	return err
}
