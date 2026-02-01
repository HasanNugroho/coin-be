package dashboard

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	dailySummaries *mongo.Collection
	transactions   *mongo.Collection
	pockets        *mongo.Collection
	userCategories *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		dailySummaries: db.Collection("daily_summaries"),
		transactions:   db.Collection("transactions"),
		pockets:        db.Collection("pockets"),
		userCategories: db.Collection("user_categories"),
	}
}

func (r *Repository) CreateDailySummary(ctx context.Context, summary *DailySummary) error {
	summary.ID = primitive.NewObjectID()
	summary.CreatedAt = time.Now()
	_, err := r.dailySummaries.InsertOne(ctx, summary)
	return err
}

func (r *Repository) GetDailySummariesByDateRange(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) ([]*DailySummary, error) {
	filter := bson.M{
		"user_id": userID,
		"date": bson.M{
			"$gte": startDate,
			"$lt":  endDate,
		},
	}

	opts := options.Find().SetSort(bson.M{"date": 1})
	cursor, err := r.dailySummaries.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var summaries []*DailySummary
	if err = cursor.All(ctx, &summaries); err != nil {
		return nil, err
	}
	return summaries, nil
}

func (r *Repository) GetHistoricalSummary(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) (float64, float64, []CategoryBreakdown, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"user_id": userID,
			"date": bson.M{
				"$gte": startDate,
				"$lt":  endDate,
			},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":           nil,
			"total_income":  bson.M{"$sum": "$total_income"},
			"total_expense": bson.M{"$sum": "$total_expense"},
			"categories":    bson.M{"$push": "$category_breakdown"},
		}}},
	}

	cursor, err := r.dailySummaries.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, 0, nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		TotalIncome  float64               `bson:"total_income"`
		TotalExpense float64               `bson:"total_expense"`
		Categories   [][]CategoryBreakdown `bson:"categories"`
	}

	if err = cursor.All(ctx, &results); err != nil {
		return 0, 0, nil, err
	}

	if len(results) == 0 {
		return 0, 0, []CategoryBreakdown{}, nil
	}

	categoryMap := make(map[string]*CategoryBreakdown)
	for _, dayCategories := range results[0].Categories {
		for _, cat := range dayCategories {
			key := cat.Type + "_"
			if cat.CategoryID != nil {
				key += cat.CategoryID.Hex()
			} else {
				key += "uncategorized"
			}

			if existing, ok := categoryMap[key]; ok {
				existing.Amount += cat.Amount
			} else {
				categoryMap[key] = &CategoryBreakdown{
					CategoryID:   cat.CategoryID,
					CategoryName: cat.CategoryName,
					Type:         cat.Type,
					Amount:       cat.Amount,
				}
			}
		}
	}

	categories := make([]CategoryBreakdown, 0, len(categoryMap))
	for _, cat := range categoryMap {
		categories = append(categories, *cat)
	}

	return results[0].TotalIncome, results[0].TotalExpense, categories, nil
}

func (r *Repository) GetLiveDeltaSummary(ctx context.Context, userID primitive.ObjectID, startDate time.Time) (float64, float64, []CategoryBreakdown, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"user_id":    userID,
			"deleted_at": nil,
			"date": bson.M{
				"$gte": startDate,
			},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id": bson.M{
				"type":        "$type",
				"category_id": "$category_id",
			},
			"amount": bson.M{"$sum": "$amount"},
		}}},
	}

	cursor, err := r.transactions.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, 0, nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		ID struct {
			Type       string              `bson:"type"`
			CategoryID *primitive.ObjectID `bson:"category_id"`
		} `bson:"_id"`
		Amount float64 `bson:"amount"`
	}

	if err = cursor.All(ctx, &results); err != nil {
		return 0, 0, nil, err
	}

	var totalIncome, totalExpense float64
	categoryMap := make(map[string]*CategoryBreakdown)

	for _, result := range results {
		txType := result.ID.Type
		amount := result.Amount

		if txType == "income" {
			totalIncome += amount
		} else if txType == "expense" {
			totalExpense += amount
		}

		if txType == "income" || txType == "expense" {
			key := txType + "_"
			if result.ID.CategoryID != nil {
				key += result.ID.CategoryID.Hex()
			} else {
				key += "uncategorized"
			}

			categoryName := "Uncategorized"
			if result.ID.CategoryID != nil {
				var category struct {
					Name string `bson:"name"`
				}
				err := r.userCategories.FindOne(ctx, bson.M{"_id": result.ID.CategoryID}).Decode(&category)
				if err == nil {
					categoryName = category.Name
				}
			}

			categoryMap[key] = &CategoryBreakdown{
				CategoryID:   result.ID.CategoryID,
				CategoryName: categoryName,
				Type:         txType,
				Amount:       amount,
			}
		}
	}

	categories := make([]CategoryBreakdown, 0, len(categoryMap))
	for _, cat := range categoryMap {
		categories = append(categories, *cat)
	}

	return totalIncome, totalExpense, categories, nil
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

func (r *Repository) GetDailyCashFlowTrend(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) ([]ChartDataPoint, error) {
	historicalSummaries, err := r.GetDailySummariesByDateRange(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	dataMap := make(map[string]*ChartDataPoint)
	for _, summary := range historicalSummaries {
		dateStr := summary.Date.Format("2006-01-02")
		dataMap[dateStr] = &ChartDataPoint{
			Date:    dateStr,
			Income:  summary.TotalIncome,
			Expense: summary.TotalExpense,
		}
	}

	todayStart := time.Now().Truncate(24 * time.Hour)
	if !startDate.After(todayStart) && todayStart.Before(endDate) {
		liveIncome, liveExpense, _, err := r.GetLiveDeltaSummary(ctx, userID, todayStart)
		if err == nil {
			dateStr := todayStart.Format("2006-01-02")
			dataMap[dateStr] = &ChartDataPoint{
				Date:    dateStr,
				Income:  liveIncome,
				Expense: liveExpense,
			}
		}
	}

	result := make([]ChartDataPoint, 0, len(dataMap))
	current := startDate
	for current.Before(endDate) {
		dateStr := current.Format("2006-01-02")
		if point, exists := dataMap[dateStr]; exists {
			result = append(result, *point)
		} else {
			result = append(result, ChartDataPoint{
				Date:    dateStr,
				Income:  0,
				Expense: 0,
			})
		}
		current = current.AddDate(0, 0, 1)
	}

	return result, nil
}

func (r *Repository) GenerateDailySummaryForDate(ctx context.Context, userID primitive.ObjectID, date time.Time) error {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.AddDate(0, 0, 1)

	filter := bson.M{
		"user_id": userID,
		"date":    startOfDay,
	}
	r.dailySummaries.DeleteOne(ctx, filter)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"user_id":    userID,
			"deleted_at": nil,
			"date": bson.M{
				"$gte": startOfDay,
				"$lt":  endOfDay,
			},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id": bson.M{
				"type":        "$type",
				"category_id": "$category_id",
			},
			"amount": bson.M{"$sum": "$amount"},
		}}},
	}

	cursor, err := r.transactions.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var results []struct {
		ID struct {
			Type       string              `bson:"type"`
			CategoryID *primitive.ObjectID `bson:"category_id"`
		} `bson:"_id"`
		Amount float64 `bson:"amount"`
	}

	if err = cursor.All(ctx, &results); err != nil {
		return err
	}

	var totalIncome, totalExpense float64
	categoryMap := make(map[string]*CategoryBreakdown)

	for _, result := range results {
		txType := result.ID.Type
		amount := result.Amount

		if txType == "income" {
			totalIncome += amount
		} else if txType == "expense" {
			totalExpense += amount
		}

		if txType == "income" || txType == "expense" {
			key := txType + "_"
			if result.ID.CategoryID != nil {
				key += result.ID.CategoryID.Hex()
			} else {
				key += "uncategorized"
			}

			categoryName := "Uncategorized"
			if result.ID.CategoryID != nil {
				var category struct {
					Name string `bson:"name"`
				}
				err := r.userCategories.FindOne(ctx, bson.M{"_id": result.ID.CategoryID}).Decode(&category)
				if err == nil {
					categoryName = category.Name
				}
			}

			categoryMap[key] = &CategoryBreakdown{
				CategoryID:   result.ID.CategoryID,
				CategoryName: categoryName,
				Type:         txType,
				Amount:       amount,
			}
		}
	}

	categoryBreakdown := make([]CategoryBreakdown, 0, len(categoryMap))
	for _, cat := range categoryMap {
		categoryBreakdown = append(categoryBreakdown, *cat)
	}

	summary := &DailySummary{
		UserID:            userID,
		Date:              startOfDay,
		TotalIncome:       totalIncome,
		TotalExpense:      totalExpense,
		CategoryBreakdown: categoryBreakdown,
	}

	return r.CreateDailySummary(ctx, summary)
}

func (r *Repository) GetAllUsersWithTransactions(ctx context.Context, date time.Time) ([]primitive.ObjectID, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.AddDate(0, 0, 1)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"deleted_at": nil,
			"date": bson.M{
				"$gte": startOfDay,
				"$lt":  endOfDay,
			},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id": "$user_id",
		}}},
	}

	cursor, err := r.transactions.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		ID primitive.ObjectID `bson:"_id"`
	}

	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	userIDs := make([]primitive.ObjectID, len(results))
	for i, result := range results {
		userIDs[i] = result.ID
	}

	return userIDs, nil
}

func (r *Repository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "date", Value: -1},
			},
			Options: options.Index().
				SetName("idx_daily_summaries_user_date").
				SetUnique(true),
		},
	}

	_, err := r.dailySummaries.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return err
	}

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

	_, err = r.transactions.Indexes().CreateMany(ctx, transactionIndexes)
	return err
}
