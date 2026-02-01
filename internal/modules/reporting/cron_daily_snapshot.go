package reporting

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DailySnapshotGenerator struct {
	db   *mongo.Database
	repo *Repository
}

func NewDailySnapshotGenerator(db *mongo.Database, repo *Repository) *DailySnapshotGenerator {
	return &DailySnapshotGenerator{
		db:   db,
		repo: repo,
	}
}

// GenerateForAllUsers generates daily snapshots for all users at 00:00 UTC
// Captures point-in-time balance state for the previous day
func (g *DailySnapshotGenerator) GenerateForAllUsers(ctx context.Context, snapshotDate time.Time) error {
	start := time.Now()
	log.Printf("[DailySnapshotGenerator] Starting snapshot generation for %s", snapshotDate.Format("2006-01-02"))

	// Get all unique user IDs
	userIDs, err := g.getAllUserIDs(ctx)
	if err != nil {
		log.Printf("[DailySnapshotGenerator] Error getting user IDs: %v", err)
		return err
	}

	log.Printf("[DailySnapshotGenerator] Found %d users to process", len(userIDs))

	// Process users in batches
	batchSize := 100
	successCount := 0
	errorCount := 0

	for i := 0; i < len(userIDs); i += batchSize {
		end := i + batchSize
		if end > len(userIDs) {
			end = len(userIDs)
		}

		batch := userIDs[i:end]
		for _, userID := range batch {
			if err := g.GenerateForUser(ctx, userID, snapshotDate); err != nil {
				log.Printf("[DailySnapshotGenerator] Error for user %v: %v", userID, err)
				errorCount++
			} else {
				successCount++
			}
		}

		log.Printf("[DailySnapshotGenerator] Batch %d-%d: %d success, %d errors", i, end, successCount, errorCount)
	}

	duration := time.Since(start)
	log.Printf("[DailySnapshotGenerator] Completed in %v: %d success, %d failed", duration, successCount, errorCount)

	return nil
}

// GenerateForUser generates a daily snapshot for a specific user
func (g *DailySnapshotGenerator) GenerateForUser(ctx context.Context, userID primitive.ObjectID, snapshotDate time.Time) error {
	// Normalize date to start of day
	snapshotDate = time.Date(snapshotDate.Year(), snapshotDate.Month(), snapshotDate.Day(), 0, 0, 0, 0, time.UTC)

	// Check if snapshot already exists
	existing, err := g.repo.GetDailySnapshot(ctx, userID, snapshotDate)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	// If complete snapshot exists, skip
	if existing != nil && existing.IsComplete {
		return nil
	}

	// Get all pockets for user
	pockets, err := g.getPocketsForUser(ctx, userID)
	if err != nil {
		return err
	}

	// Get daily report for aggregates
	dailyReport, err := g.repo.GetDailyReport(ctx, userID, snapshotDate)

	// Build snapshot
	snapshot := &DailyFinancialSnapshot{
		UserID:           userID,
		SnapshotDate:     snapshotDate,
		PocketBalances:   make([]PocketBalanceSnapshot, 0, len(pockets)),
		TotalBalance:     primitive.NewDecimal128(0, 0),
		TotalIncome:      primitive.NewDecimal128(0, 0),
		TotalExpense:     primitive.NewDecimal128(0, 0),
		YTDIncome:        primitive.NewDecimal128(0, 0),
		YTDExpense:       primitive.NewDecimal128(0, 0),
		YTDNet:           primitive.NewDecimal128(0, 0),
		TransactionCount: 0,
		IsComplete:       true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Add pocket balances
	for _, pocket := range pockets {
		snapshot.PocketBalances = append(snapshot.PocketBalances, PocketBalanceSnapshot{
			PocketID:   pocket["_id"].(primitive.ObjectID),
			PocketName: pocket["name"].(string),
			PocketType: pocket["type"].(string),
			Balance:    pocket["balance"].(primitive.Decimal128),
		})
	}

	// Add daily report aggregates if available
	if dailyReport != nil {
		snapshot.TotalIncome = dailyReport.TotalIncome
		snapshot.TotalExpense = dailyReport.TotalExpense
		snapshot.TransactionCount = int32(len(dailyReport.TransactionsByPocket))
	}

	// Upsert snapshot
	if err := g.repo.UpsertDailySnapshot(ctx, snapshot); err != nil {
		return err
	}

	return nil
}

// getAllUserIDs retrieves all unique user IDs from the system
func (g *DailySnapshotGenerator) getAllUserIDs(ctx context.Context) ([]primitive.ObjectID, error) {
	collection := g.db.Collection("pockets")

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$user_id"},
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	userIDs := make([]primitive.ObjectID, 0, len(results))
	for _, result := range results {
		if userID, ok := result["_id"].(primitive.ObjectID); ok {
			userIDs = append(userIDs, userID)
		}
	}

	return userIDs, nil
}

// getPocketsForUser retrieves all pockets for a user with current balances
func (g *DailySnapshotGenerator) getPocketsForUser(ctx context.Context, userID primitive.ObjectID) ([]bson.M, error) {
	collection := g.db.Collection("pockets")

	filter := bson.D{
		{Key: "user_id", Value: userID},
		{Key: "deleted_at", Value: bson.D{{Key: "$eq", Value: nil}}},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pockets []bson.M
	if err = cursor.All(ctx, &pockets); err != nil {
		return nil, err
	}

	return pockets, nil
}
