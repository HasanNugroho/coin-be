package reporting

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// IndexManager handles creation of all indexes for reporting collections
type IndexManager struct {
	db *mongo.Database
}

func NewIndexManager(db *mongo.Database) *IndexManager {
	return &IndexManager{db: db}
}

// CreateAllIndexes creates all necessary indexes for reporting collections
func (im *IndexManager) CreateAllIndexes(ctx context.Context) error {
	if err := im.createDailyReportIndexes(ctx); err != nil {
		return err
	}
	if err := im.createDailySnapshotIndexes(ctx); err != nil {
		return err
	}
	if err := im.createMonthlySnapshotIndexes(ctx); err != nil {
		return err
	}
	if err := im.createPocketBalanceSnapshotIndexes(ctx); err != nil {
		return err
	}
	if err := im.createAITransactionEnrichmentIndexes(ctx); err != nil {
		return err
	}
	if err := im.createAISpendingPatternIndexes(ctx); err != nil {
		return err
	}
	if err := im.createAIFinancialInsightIndexes(ctx); err != nil {
		return err
	}
	return nil
}

// ============================================================================
// DAILY FINANCIAL REPORTS INDEXES
// ============================================================================

func (im *IndexManager) createDailyReportIndexes(ctx context.Context) error {
	collection := im.db.Collection("daily_financial_reports")

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "report_date", Value: -1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "report_date", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "report_date", Value: -1}, {Key: "is_final", Value: 1}},
		},
		{
			Keys:    bson.D{{Key: "created_at", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(2592000).SetPartialFilterExpression(bson.M{"is_final": false}),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	return err
}

// ============================================================================
// DAILY FINANCIAL SNAPSHOTS INDEXES
// ============================================================================

func (im *IndexManager) createDailySnapshotIndexes(ctx context.Context) error {
	collection := im.db.Collection("daily_financial_snapshots")

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "snapshot_date", Value: -1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "snapshot_date", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "is_complete", Value: 1}, {Key: "snapshot_date", Value: -1}},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	return err
}

// ============================================================================
// MONTHLY FINANCIAL SUMMARIES INDEXES
// ============================================================================

func (im *IndexManager) createMonthlySnapshotIndexes(ctx context.Context) error {
	collection := im.db.Collection("monthly_financial_summaries")

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "month", Value: -1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "month", Value: 1}},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	return err
}

// ============================================================================
// POCKET BALANCE SNAPSHOTS INDEXES
// ============================================================================

func (im *IndexManager) createPocketBalanceSnapshotIndexes(ctx context.Context) error {
	collection := im.db.Collection("pocket_balance_snapshots")

	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "pocket_id", Value: 1}, {Key: "snapshot_time", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "snapshot_time", Value: -1}},
		},
		{
			Keys:    bson.D{{Key: "snapshot_time", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(7776000),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	return err
}

// ============================================================================
// AI TRANSACTION ENRICHMENT INDEXES
// ============================================================================

func (im *IndexManager) createAITransactionEnrichmentIndexes(ctx context.Context) error {
	collection := im.db.Collection("ai_transaction_enrichment")

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "transaction_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "date", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "is_anomaly", Value: 1}, {Key: "date", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "is_recurring", Value: 1}},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	return err
}

// ============================================================================
// AI SPENDING PATTERNS INDEXES
// ============================================================================

func (im *IndexManager) createAISpendingPatternIndexes(ctx context.Context) error {
	collection := im.db.Collection("ai_spending_patterns")

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "category", Value: 1}, {Key: "merchant", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "frequency_per_month", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "is_recurring", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "is_seasonal", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "is_essential", Value: 1}},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	return err
}

// ============================================================================
// AI FINANCIAL INSIGHTS INDEXES
// ============================================================================

func (im *IndexManager) createAIFinancialInsightIndexes(ctx context.Context) error {
	collection := im.db.Collection("ai_financial_insights")

	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "created_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "is_read", Value: 1}, {Key: "severity", Value: 1}},
		},
		{
			Keys:    bson.D{{Key: "expires_at", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(0),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	return err
}
