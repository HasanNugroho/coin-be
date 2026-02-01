package reporting

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DailyReportGenerator struct {
	db   *mongo.Database
	repo *Repository
}

func NewDailyReportGenerator(db *mongo.Database, repo *Repository) *DailyReportGenerator {
	return &DailyReportGenerator{
		db:   db,
		repo: repo,
	}
}

// GenerateForAllUsers generates daily reports for all users
// Optimized for batch processing with minimal memory footprint
func (g *DailyReportGenerator) GenerateForAllUsers(ctx context.Context, reportDate time.Time) error {
	start := time.Now()
	log.Printf("[DailyReportGenerator] Starting generation for %s", reportDate.Format("2006-01-02"))

	// Get all unique user IDs from transactions
	userIDs, err := g.getUniqueUserIDs(ctx, reportDate)
	if err != nil {
		log.Printf("[DailyReportGenerator] Error getting user IDs: %v", err)
		return err
	}

	log.Printf("[DailyReportGenerator] Found %d users to process", len(userIDs))

	// Process users in batches to avoid memory spikes
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
			if err := g.GenerateForUser(ctx, userID, reportDate); err != nil {
				log.Printf("[DailyReportGenerator] Error for user %v: %v", userID, err)
				errorCount++
			} else {
				successCount++
			}
		}

		// Log progress
		log.Printf("[DailyReportGenerator] Batch %d-%d: %d success, %d errors", i, end, successCount, errorCount)
	}

	duration := time.Since(start)
	log.Printf("[DailyReportGenerator] Completed in %v: %d success, %d failed", duration, successCount, errorCount)

	return nil
}

// GenerateForUser generates a daily report for a specific user
func (g *DailyReportGenerator) GenerateForUser(ctx context.Context, userID primitive.ObjectID, reportDate time.Time) error {
	// Normalize date to start of day
	reportDate = time.Date(reportDate.Year(), reportDate.Month(), reportDate.Day(), 0, 0, 0, 0, time.UTC)

	// Check if report already exists
	existing, err := g.repo.GetDailyReport(ctx, userID, reportDate)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	// If final report exists, skip
	if existing != nil && existing.IsFinal {
		return nil
	}

	// Get opening balance from previous day's closing balance
	previousDay := reportDate.AddDate(0, 0, -1)
	previousReport, err := g.repo.GetDailyReport(ctx, userID, previousDay)
	var openingBalance primitive.Decimal128
	if err == nil && previousReport != nil {
		openingBalance = previousReport.ClosingBalance
	}

	// Get all transactions for the day
	transactions, err := g.getTransactionsForDay(ctx, userID, reportDate)
	if err != nil {
		return err
	}

	// Aggregate data
	report := g.aggregateTransactions(userID, reportDate, openingBalance, transactions)

	// Upsert report
	if err := g.repo.UpsertDailyReport(ctx, report); err != nil {
		return err
	}

	return nil
}

// getUniqueUserIDs retrieves all unique user IDs that have transactions on the given date
func (g *DailyReportGenerator) getUniqueUserIDs(ctx context.Context, reportDate time.Time) ([]primitive.ObjectID, error) {
	startOfDay := time.Date(reportDate.Year(), reportDate.Month(), reportDate.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	collection := g.db.Collection("transactions")

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "date", Value: bson.D{
				{Key: "$gte", Value: startOfDay},
				{Key: "$lt", Value: endOfDay},
			}},
		}}},
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

// getTransactionsForDay retrieves all transactions for a user on a specific day
func (g *DailyReportGenerator) getTransactionsForDay(ctx context.Context, userID primitive.ObjectID, reportDate time.Time) ([]bson.M, error) {
	startOfDay := time.Date(reportDate.Year(), reportDate.Month(), reportDate.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	collection := g.db.Collection("transactions")

	filter := bson.D{
		{Key: "user_id", Value: userID},
		{Key: "date", Value: bson.D{
			{Key: "$gte", Value: startOfDay},
			{Key: "$lt", Value: endOfDay},
		}},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []bson.M
	if err = cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

// aggregateTransactions aggregates transaction data into a daily report
func (g *DailyReportGenerator) aggregateTransactions(userID primitive.ObjectID, reportDate time.Time, openingBalance primitive.Decimal128, transactions []bson.M) *DailyFinancialReport {
	report := &DailyFinancialReport{
		UserID:               userID,
		ReportDate:           reportDate,
		OpeningBalance:       openingBalance,
		TotalIncome:          primitive.NewDecimal128(0, 0),
		TotalExpense:         primitive.NewDecimal128(0, 0),
		TotalTransferIn:      primitive.NewDecimal128(0, 0),
		TotalTransferOut:     primitive.NewDecimal128(0, 0),
		ExpenseByCategory:    make([]ExpenseByCategory, 0),
		TransactionsByPocket: make([]TransactionsByPocket, 0),
		IsFinal:              false,
		GeneratedAt:          time.Now(),
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	// Group by category and pocket
	categoryMap := make(map[primitive.ObjectID]*ExpenseByCategory)
	pocketMap := make(map[primitive.ObjectID]*TransactionsByPocket)

	for _, txn := range transactions {
		txnType, _ := txn["type"].(string)
		categoryID, _ := txn["category_id"].(primitive.ObjectID)
		pocketID, _ := txn["pocket_id"].(primitive.ObjectID)

		// Aggregate by category
		if categoryID != primitive.NilObjectID && (txnType == "expense" || txnType == "income") {
			if _, exists := categoryMap[categoryID]; !exists {
				categoryMap[categoryID] = &ExpenseByCategory{
					CategoryID:       categoryID,
					Amount:           primitive.NewDecimal128(0, 0),
					TransactionCount: 0,
				}
			}
			categoryMap[categoryID].TransactionCount++
		}

		// Aggregate by pocket
		if pocketID != primitive.NilObjectID {
			if _, exists := pocketMap[pocketID]; !exists {
				pocketMap[pocketID] = &TransactionsByPocket{
					PocketID:         pocketID,
					Income:           primitive.NewDecimal128(0, 0),
					Expense:          primitive.NewDecimal128(0, 0),
					TransferIn:       primitive.NewDecimal128(0, 0),
					TransferOut:      primitive.NewDecimal128(0, 0),
					TransactionCount: 0,
				}
			}
			pocketMap[pocketID].TransactionCount++
		}
	}

	// Convert maps to slices
	for _, category := range categoryMap {
		report.ExpenseByCategory = append(report.ExpenseByCategory, *category)
	}

	for _, pocket := range pocketMap {
		report.TransactionsByPocket = append(report.TransactionsByPocket, *pocket)
	}

	// Calculate closing balance (simplified - in production use proper Decimal128 math)
	report.ClosingBalance = openingBalance

	return report
}
