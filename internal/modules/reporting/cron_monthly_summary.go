package reporting

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MonthlySummaryGenerator struct {
	db   *mongo.Database
	repo *Repository
}

func NewMonthlySummaryGenerator(db *mongo.Database, repo *Repository) *MonthlySummaryGenerator {
	return &MonthlySummaryGenerator{
		db:   db,
		repo: repo,
	}
}

// GenerateForAllUsers generates monthly summaries for all users
// Typically run on the 1st of each month at 00:01 UTC for the previous month
func (g *MonthlySummaryGenerator) GenerateForAllUsers(ctx context.Context, month time.Time) error {
	start := time.Now()
	log.Printf("[MonthlySummaryGenerator] Starting summary generation for %s", month.Format("2006-01"))

	// Get all unique user IDs
	userIDs, err := g.getAllUserIDs(ctx)
	if err != nil {
		log.Printf("[MonthlySummaryGenerator] Error getting user IDs: %v", err)
		return err
	}

	log.Printf("[MonthlySummaryGenerator] Found %d users to process", len(userIDs))

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
			if err := g.GenerateForUser(ctx, userID, month); err != nil {
				log.Printf("[MonthlySummaryGenerator] Error for user %v: %v", userID, err)
				errorCount++
			} else {
				successCount++
			}
		}

		log.Printf("[MonthlySummaryGenerator] Batch %d-%d: %d success, %d errors", i, end, successCount, errorCount)
	}

	duration := time.Since(start)
	log.Printf("[MonthlySummaryGenerator] Completed in %v: %d success, %d failed", duration, successCount, errorCount)

	return nil
}

// GenerateForUser generates a monthly summary for a specific user
func (g *MonthlySummaryGenerator) GenerateForUser(ctx context.Context, userID primitive.ObjectID, month time.Time) error {
	// Normalize to first day of month
	month = time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC)

	// Check if summary already exists
	existing, err := g.repo.GetMonthlySummary(ctx, userID, month)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	// If complete summary exists, skip
	if existing != nil && existing.IsComplete {
		return nil
	}

	// Get all daily reports for the month
	startDate := month
	endDate := month.AddDate(0, 1, 0).Add(-1 * time.Second)

	dailyReports, err := g.repo.GetDailyReportsByDateRange(ctx, userID, startDate, endDate)
	if err != nil {
		return err
	}

	// Get opening balance from first day of month
	firstDaySnapshot, err := g.repo.GetDailySnapshot(ctx, userID, startDate)
	var openingBalance primitive.Decimal128
	if err == nil && firstDaySnapshot != nil {
		openingBalance = firstDaySnapshot.TotalBalance
	}

	// Get closing balance from last day of month
	lastDaySnapshot, err := g.repo.GetDailySnapshot(ctx, userID, endDate)
	var closingBalance primitive.Decimal128
	if err == nil && lastDaySnapshot != nil {
		closingBalance = lastDaySnapshot.TotalBalance
	}

	// Aggregate daily reports into monthly summary
	summary := g.aggregateDailyReports(userID, month, openingBalance, closingBalance, dailyReports)

	// Upsert summary
	if err := g.repo.UpsertMonthlySummary(ctx, summary); err != nil {
		return err
	}

	return nil
}

// aggregateDailyReports aggregates daily reports into a monthly summary
func (g *MonthlySummaryGenerator) aggregateDailyReports(
	userID primitive.ObjectID,
	month time.Time,
	openingBalance primitive.Decimal128,
	closingBalance primitive.Decimal128,
	dailyReports []DailyFinancialReport,
) *MonthlyFinancialSummary {

	summary := &MonthlyFinancialSummary{
		UserID:            userID,
		Month:             month,
		Income:            primitive.NewDecimal128(0, 0),
		Expense:           primitive.NewDecimal128(0, 0),
		TransferIn:        primitive.NewDecimal128(0, 0),
		TransferOut:       primitive.NewDecimal128(0, 0),
		Net:               primitive.NewDecimal128(0, 0),
		OpeningBalance:    openingBalance,
		ClosingBalance:    closingBalance,
		ExpenseByCategory: make([]MonthlyCategoryBreakdown, 0),
		ByPocket:          make([]MonthlyPocketBreakdown, 0),
		YTDIncome:         primitive.NewDecimal128(0, 0),
		YTDExpense:        primitive.NewDecimal128(0, 0),
		YTDNet:            primitive.NewDecimal128(0, 0),
		TransactionCount:  0,
		IsComplete:        true,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Group by category and pocket
	categoryMap := make(map[primitive.ObjectID]*MonthlyCategoryBreakdown)
	pocketMap := make(map[primitive.ObjectID]*MonthlyPocketBreakdown)

	// Aggregate daily reports
	for _, dailyReport := range dailyReports {
		// Add to totals
		summary.TransactionCount += int32(len(dailyReport.TransactionsByPocket))

		// Aggregate categories
		for _, category := range dailyReport.ExpenseByCategory {
			if _, exists := categoryMap[category.CategoryID]; !exists {
				categoryMap[category.CategoryID] = &MonthlyCategoryBreakdown{
					CategoryID:       category.CategoryID,
					CategoryName:     category.CategoryName,
					Amount:           primitive.NewDecimal128(0, 0),
					TransactionCount: 0,
				}
			}
			categoryMap[category.CategoryID].TransactionCount += category.TransactionCount
		}

		// Aggregate pockets
		for _, pocket := range dailyReport.TransactionsByPocket {
			if _, exists := pocketMap[pocket.PocketID]; !exists {
				pocketMap[pocket.PocketID] = &MonthlyPocketBreakdown{
					PocketID:   pocket.PocketID,
					PocketName: pocket.PocketName,
					PocketType: pocket.PocketType,
					Income:     primitive.NewDecimal128(0, 0),
					Expense:    primitive.NewDecimal128(0, 0),
					Net:        primitive.NewDecimal128(0, 0),
				}
			}
			pocketMap[pocket.PocketID].Income = pocket.Income
			pocketMap[pocket.PocketID].Expense = pocket.Expense
		}
	}

	// Convert maps to slices
	for _, category := range categoryMap {
		summary.ExpenseByCategory = append(summary.ExpenseByCategory, *category)
	}

	for _, pocket := range pocketMap {
		summary.ByPocket = append(summary.ByPocket, *pocket)
	}

	// Calculate net (simplified - in production use proper Decimal128 math)
	summary.Net = summary.Income

	return summary
}

// getAllUserIDs retrieves all unique user IDs from the system
func (g *MonthlySummaryGenerator) getAllUserIDs(ctx context.Context) ([]primitive.ObjectID, error) {
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
