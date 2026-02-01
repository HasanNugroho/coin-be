package reporting

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{db: db}
}

// ============================================================================
// DAILY FINANCIAL REPORTS
// ============================================================================

func (r *Repository) CreateDailyReport(ctx context.Context, report *DailyFinancialReport) error {
	report.CreatedAt = time.Now()
	report.UpdatedAt = time.Now()

	_, err := r.db.Collection("daily_financial_reports").InsertOne(ctx, report)
	return err
}

func (r *Repository) UpsertDailyReport(ctx context.Context, report *DailyFinancialReport) error {
	report.UpdatedAt = time.Now()
	if report.CreatedAt.IsZero() {
		report.CreatedAt = time.Now()
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.db.Collection("daily_financial_reports").UpdateOne(
		ctx,
		bson.M{
			"user_id":     report.UserID,
			"report_date": report.ReportDate,
		},
		bson.M{"$set": report},
		opts,
	)
	return err
}

func (r *Repository) GetDailyReport(ctx context.Context, userID primitive.ObjectID, reportDate time.Time) (*DailyFinancialReport, error) {
	var report DailyFinancialReport
	err := r.db.Collection("daily_financial_reports").FindOne(ctx, bson.M{
		"user_id":     userID,
		"report_date": reportDate,
	}).Decode(&report)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &report, err
}

func (r *Repository) GetDailyReportsByDateRange(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) ([]DailyFinancialReport, error) {
	opts := options.Find().SetSort(bson.M{"report_date": -1})
	cursor, err := r.db.Collection("daily_financial_reports").Find(ctx, bson.M{
		"user_id": userID,
		"report_date": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
		"is_final": true,
	}, opts)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reports []DailyFinancialReport
	return reports, cursor.All(ctx, &reports)
}

func (r *Repository) GetLatestDailyReport(ctx context.Context, userID primitive.ObjectID) (*DailyFinancialReport, error) {
	opts := options.FindOne().SetSort(bson.M{"report_date": -1})
	var report DailyFinancialReport
	err := r.db.Collection("daily_financial_reports").FindOne(ctx, bson.M{
		"user_id":   userID,
		"is_final":  true,
		"deleted_at": nil,
	}, opts).Decode(&report)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &report, err
}

// ============================================================================
// DAILY FINANCIAL SNAPSHOTS
// ============================================================================

func (r *Repository) CreateDailySnapshot(ctx context.Context, snapshot *DailyFinancialSnapshot) error {
	snapshot.CreatedAt = time.Now()
	snapshot.UpdatedAt = time.Now()

	_, err := r.db.Collection("daily_financial_snapshots").InsertOne(ctx, snapshot)
	return err
}

func (r *Repository) UpsertDailySnapshot(ctx context.Context, snapshot *DailyFinancialSnapshot) error {
	snapshot.UpdatedAt = time.Now()
	if snapshot.CreatedAt.IsZero() {
		snapshot.CreatedAt = time.Now()
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.db.Collection("daily_financial_snapshots").UpdateOne(
		ctx,
		bson.M{
			"user_id":       snapshot.UserID,
			"snapshot_date": snapshot.SnapshotDate,
		},
		bson.M{"$set": snapshot},
		opts,
	)
	return err
}

func (r *Repository) GetDailySnapshot(ctx context.Context, userID primitive.ObjectID, snapshotDate time.Time) (*DailyFinancialSnapshot, error) {
	var snapshot DailyFinancialSnapshot
	err := r.db.Collection("daily_financial_snapshots").FindOne(ctx, bson.M{
		"user_id":       userID,
		"snapshot_date": snapshotDate,
	}).Decode(&snapshot)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &snapshot, err
}

func (r *Repository) GetDailySnapshotsByDateRange(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) ([]DailyFinancialSnapshot, error) {
	opts := options.Find().SetSort(bson.M{"snapshot_date": -1})
	cursor, err := r.db.Collection("daily_financial_snapshots").Find(ctx, bson.M{
		"user_id": userID,
		"snapshot_date": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
		"is_complete": true,
	}, opts)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var snapshots []DailyFinancialSnapshot
	return snapshots, cursor.All(ctx, &snapshots)
}

// ============================================================================
// MONTHLY FINANCIAL SUMMARIES
// ============================================================================

func (r *Repository) CreateMonthlySummary(ctx context.Context, summary *MonthlyFinancialSummary) error {
	summary.CreatedAt = time.Now()
	summary.UpdatedAt = time.Now()

	_, err := r.db.Collection("monthly_financial_summaries").InsertOne(ctx, summary)
	return err
}

func (r *Repository) UpsertMonthlySummary(ctx context.Context, summary *MonthlyFinancialSummary) error {
	summary.UpdatedAt = time.Now()
	if summary.CreatedAt.IsZero() {
		summary.CreatedAt = time.Now()
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.db.Collection("monthly_financial_summaries").UpdateOne(
		ctx,
		bson.M{
			"user_id": summary.UserID,
			"month":   summary.Month,
		},
		bson.M{"$set": summary},
		opts,
	)
	return err
}

func (r *Repository) GetMonthlySummary(ctx context.Context, userID primitive.ObjectID, month time.Time) (*MonthlyFinancialSummary, error) {
	var summary MonthlyFinancialSummary
	err := r.db.Collection("monthly_financial_summaries").FindOne(ctx, bson.M{
		"user_id": userID,
		"month":   month,
	}).Decode(&summary)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &summary, err
}

func (r *Repository) GetMonthlySummariesByRange(ctx context.Context, userID primitive.ObjectID, startMonth, endMonth time.Time) ([]MonthlyFinancialSummary, error) {
	opts := options.Find().SetSort(bson.M{"month": -1})
	cursor, err := r.db.Collection("monthly_financial_summaries").Find(ctx, bson.M{
		"user_id": userID,
		"month": bson.M{
			"$gte": startMonth,
			"$lte": endMonth,
		},
	}, opts)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var summaries []MonthlyFinancialSummary
	return summaries, cursor.All(ctx, &summaries)
}

// ============================================================================
// POCKET BALANCE HISTORY SNAPSHOTS
// ============================================================================

func (r *Repository) CreateBalanceSnapshot(ctx context.Context, snapshot *PocketBalanceHistorySnapshot) error {
	snapshot.CreatedAt = time.Now()

	_, err := r.db.Collection("pocket_balance_snapshots").InsertOne(ctx, snapshot)
	return err
}

func (r *Repository) GetBalanceSnapshotsByPocket(ctx context.Context, userID, pocketID primitive.ObjectID, limit int64) ([]PocketBalanceHistorySnapshot, error) {
	opts := options.Find().SetSort(bson.M{"snapshot_time": -1}).SetLimit(limit)
	cursor, err := r.db.Collection("pocket_balance_snapshots").Find(ctx, bson.M{
		"user_id":   userID,
		"pocket_id": pocketID,
	}, opts)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var snapshots []PocketBalanceHistorySnapshot
	return snapshots, cursor.All(ctx, &snapshots)
}

func (r *Repository) GetBalanceSnapshotsByDateRange(ctx context.Context, userID primitive.ObjectID, startTime, endTime time.Time) ([]PocketBalanceHistorySnapshot, error) {
	opts := options.Find().SetSort(bson.M{"snapshot_time": -1})
	cursor, err := r.db.Collection("pocket_balance_snapshots").Find(ctx, bson.M{
		"user_id": userID,
		"snapshot_time": bson.M{
			"$gte": startTime,
			"$lte": endTime,
		},
	}, opts)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var snapshots []PocketBalanceHistorySnapshot
	return snapshots, cursor.All(ctx, &snapshots)
}

// ============================================================================
// AI TRANSACTION ENRICHMENT
// ============================================================================

func (r *Repository) CreateTransactionEnrichment(ctx context.Context, enrichment *AITransactionEnrichment) error {
	enrichment.CreatedAt = time.Now()
	enrichment.UpdatedAt = time.Now()

	_, err := r.db.Collection("ai_transaction_enrichment").InsertOne(ctx, enrichment)
	return err
}

func (r *Repository) UpsertTransactionEnrichment(ctx context.Context, enrichment *AITransactionEnrichment) error {
	enrichment.UpdatedAt = time.Now()
	if enrichment.CreatedAt.IsZero() {
		enrichment.CreatedAt = time.Now()
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.db.Collection("ai_transaction_enrichment").UpdateOne(
		ctx,
		bson.M{
			"user_id":        enrichment.UserID,
			"transaction_id": enrichment.TransactionID,
		},
		bson.M{"$set": enrichment},
		opts,
	)
	return err
}

func (r *Repository) GetTransactionEnrichment(ctx context.Context, userID, transactionID primitive.ObjectID) (*AITransactionEnrichment, error) {
	var enrichment AITransactionEnrichment
	err := r.db.Collection("ai_transaction_enrichment").FindOne(ctx, bson.M{
		"user_id":        userID,
		"transaction_id": transactionID,
	}).Decode(&enrichment)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &enrichment, err
}

func (r *Repository) GetAnomalousTransactions(ctx context.Context, userID primitive.ObjectID, limit int64) ([]AITransactionEnrichment, error) {
	opts := options.Find().SetSort(bson.M{"date": -1}).SetLimit(limit)
	cursor, err := r.db.Collection("ai_transaction_enrichment").Find(ctx, bson.M{
		"user_id":     userID,
		"is_anomaly":  true,
		"deleted_at":  nil,
	}, opts)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var enrichments []AITransactionEnrichment
	return enrichments, cursor.All(ctx, &enrichments)
}

func (r *Repository) GetTransactionEnrichmentsByDateRange(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) ([]AITransactionEnrichment, error) {
	opts := options.Find().SetSort(bson.M{"date": -1})
	cursor, err := r.db.Collection("ai_transaction_enrichment").Find(ctx, bson.M{
		"user_id": userID,
		"date": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
		"deleted_at": nil,
	}, opts)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var enrichments []AITransactionEnrichment
	return enrichments, cursor.All(ctx, &enrichments)
}

// ============================================================================
// AI SPENDING PATTERNS
// ============================================================================

func (r *Repository) CreateSpendingPattern(ctx context.Context, pattern *AISpendingPattern) error {
	pattern.CreatedAt = time.Now()
	pattern.UpdatedAt = time.Now()

	_, err := r.db.Collection("ai_spending_patterns").InsertOne(ctx, pattern)
	return err
}

func (r *Repository) UpsertSpendingPattern(ctx context.Context, pattern *AISpendingPattern) error {
	pattern.UpdatedAt = time.Now()
	if pattern.CreatedAt.IsZero() {
		pattern.CreatedAt = time.Now()
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.db.Collection("ai_spending_patterns").UpdateOne(
		ctx,
		bson.M{
			"user_id":   pattern.UserID,
			"category":  pattern.Category,
			"merchant":  pattern.Merchant,
		},
		bson.M{"$set": pattern},
		opts,
	)
	return err
}

func (r *Repository) GetSpendingPattern(ctx context.Context, userID primitive.ObjectID, category, merchant string) (*AISpendingPattern, error) {
	var pattern AISpendingPattern
	err := r.db.Collection("ai_spending_patterns").FindOne(ctx, bson.M{
		"user_id":  userID,
		"category": category,
		"merchant": merchant,
	}).Decode(&pattern)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &pattern, err
}

func (r *Repository) GetTopSpendingPatterns(ctx context.Context, userID primitive.ObjectID, limit int64) ([]AISpendingPattern, error) {
	opts := options.Find().SetSort(bson.M{"frequency_per_month": -1}).SetLimit(limit)
	cursor, err := r.db.Collection("ai_spending_patterns").Find(ctx, bson.M{
		"user_id":   userID,
		"deleted_at": nil,
	}, opts)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var patterns []AISpendingPattern
	return patterns, cursor.All(ctx, &patterns)
}

func (r *Repository) GetRecurringSpendingPatterns(ctx context.Context, userID primitive.ObjectID) ([]AISpendingPattern, error) {
	cursor, err := r.db.Collection("ai_spending_patterns").Find(ctx, bson.M{
		"user_id":      userID,
		"is_recurring": true,
		"deleted_at":   nil,
	})

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var patterns []AISpendingPattern
	return patterns, cursor.All(ctx, &patterns)
}

// ============================================================================
// AI FINANCIAL INSIGHTS
// ============================================================================

func (r *Repository) CreateFinancialInsight(ctx context.Context, insight *AIFinancialInsight) error {
	insight.CreatedAt = time.Now()

	_, err := r.db.Collection("ai_financial_insights").InsertOne(ctx, insight)
	return err
}

func (r *Repository) GetFinancialInsights(ctx context.Context, userID primitive.ObjectID, limit int64) ([]AIFinancialInsight, error) {
	opts := options.Find().SetSort(bson.M{"created_at": -1}).SetLimit(limit)
	cursor, err := r.db.Collection("ai_financial_insights").Find(ctx, bson.M{
		"user_id":   userID,
		"deleted_at": nil,
	}, opts)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var insights []AIFinancialInsight
	return insights, cursor.All(ctx, &insights)
}

func (r *Repository) GetUnreadInsights(ctx context.Context, userID primitive.ObjectID) ([]AIFinancialInsight, error) {
	cursor, err := r.db.Collection("ai_financial_insights").Find(ctx, bson.M{
		"user_id":    userID,
		"is_read":    false,
		"deleted_at": nil,
	})

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var insights []AIFinancialInsight
	return insights, cursor.All(ctx, &insights)
}

func (r *Repository) GetInsightsBySeverity(ctx context.Context, userID primitive.ObjectID, severity string) ([]AIFinancialInsight, error) {
	cursor, err := r.db.Collection("ai_financial_insights").Find(ctx, bson.M{
		"user_id":    userID,
		"severity":   severity,
		"deleted_at": nil,
	})

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var insights []AIFinancialInsight
	return insights, cursor.All(ctx, &insights)
}

func (r *Repository) MarkInsightAsRead(ctx context.Context, insightID primitive.ObjectID) error {
	_, err := r.db.Collection("ai_financial_insights").UpdateOne(
		ctx,
		bson.M{"_id": insightID},
		bson.M{
			"$set": bson.M{
				"is_read": true,
				"read_at": time.Now(),
			},
		},
	)
	return err
}

func (r *Repository) DeleteInsight(ctx context.Context, insightID primitive.ObjectID) error {
	_, err := r.db.Collection("ai_financial_insights").UpdateOne(
		ctx,
		bson.M{"_id": insightID},
		bson.M{
			"$set": bson.M{
				"deleted_at": time.Now(),
			},
		},
	)
	return err
}
