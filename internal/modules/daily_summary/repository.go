package daily_summary

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
	userCategories *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		dailySummaries: db.Collection("daily_summaries"),
		transactions:   db.Collection("transactions"),
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

func (r *Repository) DeleteDailySummariesByDateRange(ctx context.Context, startDate time.Time) error {
	filter := bson.M{
		"date": bson.M{
			"$gte": startDate,
		},
	}
	_, err := r.dailySummaries.DeleteMany(ctx, filter)
	return err
}

func (r *Repository) GetHistoricalSummary(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) (float64, float64, []CategoryBreakdown, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"user_id": userID,
			"date": bson.M{
				"$gte": startDate.UTC(),
				"$lt":  endDate.UTC(),
			},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":           nil,
			"total_income":  bson.M{"$sum": "$total_income"},
			"total_expense": bson.M{"$sum": "$total_expense"},
			"categories":    bson.M{"$push": "$category_breakdown"},
		}}},
		{{Key: "$project", Value: bson.M{
			"_id":           0,
			"total_income":  1,
			"total_expense": 1,
			"categories": bson.M{
				"$reduce": bson.M{
					"input":        "$categories",
					"initialValue": bson.A{},
					"in": bson.M{
						"$concatArrays": bson.A{"$$value", "$$this"},
					},
				},
			},
		}}},
	}

	cursor, err := r.dailySummaries.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, 0, nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		TotalIncome  float64             `bson:"total_income"`
		TotalExpense float64             `bson:"total_expense"`
		Categories   []CategoryBreakdown `bson:"categories"`
	}

	if err = cursor.All(ctx, &results); err != nil {
		return 0, 0, nil, err
	}

	if len(results) == 0 {
		return 0, 0, []CategoryBreakdown{}, nil
	}

	// Dedup & merge category dengan key yang sama
	categoryMap := make(map[string]*CategoryBreakdown)
	for _, cat := range results[0].Categories {
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

	categories := make([]CategoryBreakdown, 0, len(categoryMap))
	for _, cat := range categoryMap {
		categories = append(categories, *cat)
	}

	return results[0].TotalIncome, results[0].TotalExpense, categories, nil
}

// GenerateDailySummariesFromTo menghitung daily summary dari startDate sampai kemarin (cutoff) untuk SEMUA user.
func (r *Repository) GenerateDailySummariesFromTo(ctx context.Context, startDate time.Time) error {
	loc := time.UTC
	start := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, loc)
	today := time.Now().UTC()
	cutoff := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, loc)

	if !start.Before(cutoff) {
		return nil
	}

	// 1. Find all users with transactions in the range
	userIDsInterface, err := r.transactions.Distinct(ctx, "user_id", bson.M{
		"deleted_at": nil,
		"date":       bson.M{"$gte": start, "$lt": cutoff},
	})
	if err != nil {
		return err
	}

	if len(userIDsInterface) == 0 {
		return nil
	}

	userIDs := make([]primitive.ObjectID, len(userIDsInterface))
	for i, v := range userIDsInterface {
		userIDs[i] = v.(primitive.ObjectID)
	}

	// 2. Batch process users
	batchSize := 100
	for i := 0; i < len(userIDs); i += batchSize {
		end := i + batchSize
		if end > len(userIDs) {
			end = len(userIDs)
		}
		batchUserIDs := userIDs[i:end]

		// 3. Fetch transactions per batch
		cursor, err := r.transactions.Find(ctx, bson.M{
			"user_id":    bson.M{"$in": batchUserIDs},
			"deleted_at": nil,
			"date":       bson.M{"$gte": start, "$lt": cutoff},
		})
		if err != nil {
			return err
		}

		var txs []struct {
			UserID             primitive.ObjectID  `bson:"user_id"`
			Type               string              `bson:"type"`
			Amount             float64             `bson:"amount"`
			Date               time.Time           `bson:"date"`
			CategoryID         *primitive.ObjectID `bson:"category_id"`
			PocketFromID       *primitive.ObjectID `bson:"pocket_from_id"`
			PocketToID         *primitive.ObjectID `bson:"pocket_to_id"`
			UserPlatformFromID *primitive.ObjectID `bson:"user_platform_from_id"`
			UserPlatformToID   *primitive.ObjectID `bson:"user_platform_to_id"`
		}
		if err = cursor.All(ctx, &txs); err != nil {
			cursor.Close(ctx)
			return err
		}
		cursor.Close(ctx)

		// 4. Calculate per user per day in Go
		type summaryKey struct {
			UserID primitive.ObjectID
			Day    time.Time
		}
		type dayData struct {
			TotalIncome  float64
			TotalExpense float64
			Categories   map[string]*CategoryBreakdown
			Pockets      map[string]*PocketBreakdown
			Platforms    map[string]*PlatformBreakdown
		}

		resultsMap := make(map[summaryKey]*dayData)

		for _, tx := range txs {
			if tx.Type != "income" && tx.Type != "expense" {
				continue
			}

			day := time.Date(tx.Date.Year(), tx.Date.Month(), tx.Date.Day(), 0, 0, 0, 0, loc)
			key := summaryKey{UserID: tx.UserID, Day: day}

			d, ok := resultsMap[key]
			if !ok {
				d = &dayData{
					Categories: make(map[string]*CategoryBreakdown),
					Pockets:    make(map[string]*PocketBreakdown),
					Platforms:  make(map[string]*PlatformBreakdown),
				}
				resultsMap[key] = d
			}

			if tx.Type == "income" {
				d.TotalIncome += tx.Amount
			} else {
				d.TotalExpense += tx.Amount
			}

			// Category
			ck := tx.Type + "_uncategorized"
			var catID *primitive.ObjectID
			if tx.CategoryID != nil {
				ck = tx.Type + "_" + tx.CategoryID.Hex()
				catID = tx.CategoryID
			}
			if existing, ok := d.Categories[ck]; ok {
				existing.Amount += tx.Amount
			} else {
				d.Categories[ck] = &CategoryBreakdown{
					Type:       tx.Type,
					CategoryID: catID,
					Amount:     tx.Amount,
				}
			}

			// Pocket
			var pocketID *primitive.ObjectID
			if tx.Type == "income" {
				pocketID = tx.PocketToID
			} else {
				pocketID = tx.PocketFromID
			}
			if pocketID != nil {
				pk := tx.Type + "_" + pocketID.Hex()
				if existing, ok := d.Pockets[pk]; ok {
					existing.Amount += tx.Amount
				} else {
					d.Pockets[pk] = &PocketBreakdown{
						Type:     tx.Type,
						PocketID: pocketID,
						Amount:   tx.Amount,
					}
				}
			}

			// Platform
			var platformID *primitive.ObjectID
			if tx.Type == "income" {
				platformID = tx.UserPlatformToID
			} else {
				platformID = tx.UserPlatformFromID
			}
			if platformID != nil {
				plk := tx.Type + "_" + platformID.Hex()
				if existing, ok := d.Platforms[plk]; ok {
					existing.Amount += tx.Amount
				} else {
					d.Platforms[plk] = &PlatformBreakdown{
						Type:       tx.Type,
						PlatformID: platformID,
						Amount:     tx.Amount,
					}
				}
			}
		}

		// // 5. Bulk delete
		// _, err = r.dailySummaries.DeleteMany(ctx, bson.M{
		// 	"user_id": bson.M{"$in": batchUserIDs},
		// 	"date":    bson.M{"$gte": start, "$lt": cutoff},
		// })
		// if err != nil {
		// 	return err
		// }

		// 6. Bulk insert
		if len(resultsMap) > 0 {
			now := time.Now()
			docs := make([]interface{}, 0, len(resultsMap))
			for k, v := range resultsMap {
				cats := make([]CategoryBreakdown, 0, len(v.Categories))
				for _, c := range v.Categories {
					cats = append(cats, *c)
				}
				pks := make([]PocketBreakdown, 0, len(v.Pockets))
				for _, p := range v.Pockets {
					pks = append(pks, *p)
				}
				pls := make([]PlatformBreakdown, 0, len(v.Platforms))
				for _, pl := range v.Platforms {
					pls = append(pls, *pl)
				}

				docs = append(docs, &DailySummary{
					UserID:            k.UserID,
					Date:              k.Day,
					TotalIncome:       v.TotalIncome,
					TotalExpense:      v.TotalExpense,
					CategoryBreakdown: cats,
					PocketBreakdown:   pks,
					PlatformBreakdown: pls,
					CreatedAt:         now,
				})
			}
			_, err = r.dailySummaries.InsertMany(ctx, docs)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *Repository) GenerateDailySummaryForDate(ctx context.Context, userID primitive.ObjectID, date time.Time) error {
	loc := time.UTC
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc)
	endOfDay := startOfDay.AddDate(0, 0, 1)

	// 1. Fetch transactions for the day
	cursor, err := r.transactions.Find(ctx, bson.M{
		"user_id":    userID,
		"deleted_at": nil,
		"date":       bson.M{"$gte": startOfDay, "$lt": endOfDay},
	})
	if err != nil {
		return err
	}

	var txs []struct {
		Type               string              `bson:"type"`
		Amount             float64             `bson:"amount"`
		CategoryID         *primitive.ObjectID `bson:"category_id"`
		PocketFromID       *primitive.ObjectID `bson:"pocket_from_id"`
		PocketToID         *primitive.ObjectID `bson:"pocket_to_id"`
		UserPlatformFromID *primitive.ObjectID `bson:"user_platform_from_id"`
		UserPlatformToID   *primitive.ObjectID `bson:"user_platform_to_id"`
	}
	if err = cursor.All(ctx, &txs); err != nil {
		cursor.Close(ctx)
		return err
	}
	cursor.Close(ctx)

	var totalIncome, totalExpense float64
	categoryMap := make(map[string]*CategoryBreakdown)
	pocketMap := make(map[string]*PocketBreakdown)
	platformMap := make(map[string]*PlatformBreakdown)

	for _, tx := range txs {
		if tx.Type != "income" && tx.Type != "expense" {
			continue
		}

		if tx.Type == "income" {
			totalIncome += tx.Amount
		} else {
			totalExpense += tx.Amount
		}

		// Category
		ck := tx.Type + "_uncategorized"
		var catID *primitive.ObjectID
		if tx.CategoryID != nil {
			ck = tx.Type + "_" + tx.CategoryID.Hex()
			catID = tx.CategoryID
		}
		if existing, ok := categoryMap[ck]; ok {
			existing.Amount += tx.Amount
		} else {
			categoryMap[ck] = &CategoryBreakdown{
				Type:       tx.Type,
				CategoryID: catID,
				Amount:     tx.Amount,
			}
		}

		// Pocket
		var pocketID *primitive.ObjectID
		if tx.Type == "income" {
			pocketID = tx.PocketToID
		} else {
			pocketID = tx.PocketFromID
		}
		if pocketID != nil {
			pk := tx.Type + "_" + pocketID.Hex()
			if existing, ok := pocketMap[pk]; ok {
				existing.Amount += tx.Amount
			} else {
				pocketMap[pk] = &PocketBreakdown{
					Type:     tx.Type,
					PocketID: pocketID,
					Amount:   tx.Amount,
				}
			}
		}

		// Platform
		var platformID *primitive.ObjectID
		if tx.Type == "income" {
			platformID = tx.UserPlatformToID
		} else {
			platformID = tx.UserPlatformFromID
		}
		if platformID != nil {
			plk := tx.Type + "_" + platformID.Hex()
			if existing, ok := platformMap[plk]; ok {
				existing.Amount += tx.Amount
			} else {
				platformMap[plk] = &PlatformBreakdown{
					Type:       tx.Type,
					PlatformID: platformID,
					Amount:     tx.Amount,
				}
			}
		}
	}

	categoryBreakdown := make([]CategoryBreakdown, 0, len(categoryMap))
	for _, v := range categoryMap {
		categoryBreakdown = append(categoryBreakdown, *v)
	}

	pocketBreakdown := make([]PocketBreakdown, 0, len(pocketMap))
	for _, v := range pocketMap {
		pocketBreakdown = append(pocketBreakdown, *v)
	}

	platformBreakdown := make([]PlatformBreakdown, 0, len(platformMap))
	for _, v := range platformMap {
		platformBreakdown = append(platformBreakdown, *v)
	}

	// 2. Delete existing
	_, err = r.dailySummaries.DeleteOne(ctx, bson.M{
		"user_id": userID,
		"date":    startOfDay,
	})
	if err != nil {
		return err
	}

	// 3. Insert new
	summary := &DailySummary{
		UserID:            userID,
		Date:              startOfDay,
		TotalIncome:       totalIncome,
		TotalExpense:      totalExpense,
		CategoryBreakdown: categoryBreakdown,
		PocketBreakdown:   pocketBreakdown,
		PlatformBreakdown: platformBreakdown,
		CreatedAt:         time.Now(),
	}

	_, err = r.dailySummaries.InsertOne(ctx, summary)
	return err
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
	return err
}
