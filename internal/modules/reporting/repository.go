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

func (r *Repository) CreateIndexes(ctx context.Context) error {
	collections := map[string][][]bson.E{
		"daily_financial_reports": {
			{{Key: "user_id", Value: 1}, {Key: "report_date", Value: -1}},
			{{Key: "user_id", Value: 1}, {Key: "report_date", Value: 1}},
			{{Key: "user_id", Value: 1}, {Key: "is_final", Value: 1}, {Key: "report_date", Value: -1}},
		},
		"daily_financial_snapshots": {
			{{Key: "user_id", Value: 1}, {Key: "snapshot_date", Value: -1}},
		},
		"monthly_financial_summaries": {
			{{Key: "user_id", Value: 1}, {Key: "year_month", Value: -1}},
			{{Key: "user_id", Value: 1}, {Key: "year_month", Value: 1}},
			{{Key: "user_id", Value: 1}, {Key: "is_final", Value: 1}},
		},
		"pocket_balance_snapshots": {
			{{Key: "user_id", Value: 1}, {Key: "pocket_id", Value: 1}, {Key: "snapshot_date", Value: -1}},
			{{Key: "pocket_id", Value: 1}, {Key: "snapshot_date", Value: -1}},
		},
		"ai_financial_context": {
			{{Key: "user_id", Value: 1}},
			{{Key: "user_id", Value: 1}, {Key: "updated_at", Value: -1}},
		},
	}

	for collName, indexes := range collections {
		coll := r.db.Collection(collName)
		for _, index := range indexes {
			indexModel := mongo.IndexModel{
				Keys: bson.D(index),
			}
			_, err := coll.Indexes().CreateOne(ctx, indexModel)
			if err != nil {
				return err
			}
		}
	}

	// Create unique indexes
	uniqueIndexes := map[string]bson.D{
		"daily_financial_reports": {{Key: "user_id", Value: 1}, {Key: "report_date", Value: 1}},
		"daily_financial_snapshots": {{Key: "user_id", Value: 1}, {Key: "snapshot_date", Value: 1}},
		"monthly_financial_summaries": {{Key: "user_id", Value: 1}, {Key: "year_month", Value: 1}},
	}

	for collName, indexKeys := range uniqueIndexes {
		coll := r.db.Collection(collName)
		indexModel := mongo.IndexModel{
			Keys: indexKeys,
			Options: options.Index().SetUnique(true),
		}
		_, err := coll.Indexes().CreateOne(ctx, indexModel)
		if err != nil {
			return err
		}
	}

	return nil
}

// Daily Report operations
func (r *Repository) UpsertDailyReport(ctx context.Context, report *DailyFinancialReport) error {
	coll := r.db.Collection("daily_financial_reports")
	
	report.UpdatedAt = time.Now()
	if report.CreatedAt.IsZero() {
		report.CreatedAt = time.Now()
	}

	opts := options.Update().SetUpsert(true)
	filter := bson.M{
		"user_id": report.UserID,
		"report_date": bson.M{
			"$gte": time.Date(report.ReportDate.Year(), report.ReportDate.Month(), report.ReportDate.Day(), 0, 0, 0, 0, time.UTC),
			"$lt":  time.Date(report.ReportDate.Year(), report.ReportDate.Month(), report.ReportDate.Day()+1, 0, 0, 0, 0, time.UTC),
		},
	}

	_, err := coll.UpdateOne(ctx, filter, bson.M{"$set": report}, opts)
	return err
}

func (r *Repository) GetDailyReport(ctx context.Context, userID primitive.ObjectID, date time.Time) (*DailyFinancialReport, error) {
	coll := r.db.Collection("daily_financial_reports")
	
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.AddDate(0, 0, 1)

	filter := bson.M{
		"user_id": userID,
		"report_date": bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		},
	}

	var report DailyFinancialReport
	err := coll.FindOne(ctx, filter).Decode(&report)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &report, nil
}

// Daily Snapshot operations
func (r *Repository) UpsertDailySnapshot(ctx context.Context, snapshot *DailyFinancialSnapshot) error {
	coll := r.db.Collection("daily_financial_snapshots")
	
	snapshot.UpdatedAt = time.Now()
	if snapshot.CreatedAt.IsZero() {
		snapshot.CreatedAt = time.Now()
	}

	opts := options.Update().SetUpsert(true)
	startOfDay := time.Date(snapshot.SnapshotDate.Year(), snapshot.SnapshotDate.Month(), snapshot.SnapshotDate.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.AddDate(0, 0, 1)

	filter := bson.M{
		"user_id": snapshot.UserID,
		"snapshot_date": bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		},
	}

	_, err := coll.UpdateOne(ctx, filter, bson.M{"$set": snapshot}, opts)
	return err
}

func (r *Repository) GetDailySnapshot(ctx context.Context, userID primitive.ObjectID, date time.Time) (*DailyFinancialSnapshot, error) {
	coll := r.db.Collection("daily_financial_snapshots")
	
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.AddDate(0, 0, 1)

	filter := bson.M{
		"user_id": userID,
		"snapshot_date": bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		},
	}

	var snapshot DailyFinancialSnapshot
	err := coll.FindOne(ctx, filter).Decode(&snapshot)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &snapshot, nil
}

// Monthly Summary operations
func (r *Repository) UpsertMonthlySummary(ctx context.Context, summary *MonthlyFinancialSummary) error {
	coll := r.db.Collection("monthly_financial_summaries")
	
	summary.UpdatedAt = time.Now()
	if summary.CreatedAt.IsZero() {
		summary.CreatedAt = time.Now()
	}

	opts := options.Update().SetUpsert(true)
	filter := bson.M{
		"user_id": summary.UserID,
		"year_month": summary.YearMonth,
	}

	_, err := coll.UpdateOne(ctx, filter, bson.M{"$set": summary}, opts)
	return err
}

func (r *Repository) GetMonthlySummary(ctx context.Context, userID primitive.ObjectID, yearMonth string) (*MonthlyFinancialSummary, error) {
	coll := r.db.Collection("monthly_financial_summaries")
	
	filter := bson.M{
		"user_id": userID,
		"year_month": yearMonth,
	}

	var summary MonthlyFinancialSummary
	err := coll.FindOne(ctx, filter).Decode(&summary)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &summary, nil
}

func (r *Repository) GetMonthlySummariesRange(ctx context.Context, userID primitive.ObjectID, startMonth, endMonth string, limit int64) ([]MonthlyFinancialSummary, error) {
	coll := r.db.Collection("monthly_financial_summaries")
	
	filter := bson.M{
		"user_id": userID,
		"year_month": bson.M{
			"$gte": startMonth,
			"$lte": endMonth,
		},
	}

	opts := options.Find().SetSort(bson.M{"year_month": 1}).SetLimit(limit)
	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var summaries []MonthlyFinancialSummary
	err = cursor.All(ctx, &summaries)
	return summaries, err
}

// Pocket Balance Snapshot operations
func (r *Repository) CreatePocketBalanceSnapshot(ctx context.Context, snapshot *PocketBalanceSnapshot) error {
	coll := r.db.Collection("pocket_balance_snapshots")
	
	if snapshot.CreatedAt.IsZero() {
		snapshot.CreatedAt = time.Now()
	}

	_, err := coll.InsertOne(ctx, snapshot)
	return err
}

func (r *Repository) GetPocketBalanceHistory(ctx context.Context, pocketID primitive.ObjectID, days int64) ([]PocketBalanceSnapshot, error) {
	coll := r.db.Collection("pocket_balance_snapshots")
	
	startDate := time.Now().AddDate(0, 0, -int(days))
	filter := bson.M{
		"pocket_id": pocketID,
		"snapshot_date": bson.M{
			"$gte": startDate,
		},
	}

	opts := options.Find().SetSort(bson.M{"snapshot_date": -1})
	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var snapshots []PocketBalanceSnapshot
	err = cursor.All(ctx, &snapshots)
	return snapshots, err
}

// AI Financial Context operations
func (r *Repository) UpsertAIFinancialContext(ctx context.Context, context *AIFinancialContext) error {
	coll := r.db.Collection("ai_financial_context")
	
	context.UpdatedAt = time.Now()

	opts := options.Update().SetUpsert(true)
	filter := bson.M{
		"user_id": context.UserID,
	}

	_, err := coll.UpdateOne(ctx, filter, bson.M{"$set": context}, opts)
	return err
}

func (r *Repository) GetAIFinancialContext(ctx context.Context, userID primitive.ObjectID) (*AIFinancialContext, error) {
	coll := r.db.Collection("ai_financial_context")
	
	filter := bson.M{
		"user_id": userID,
	}

	var context AIFinancialContext
	err := coll.FindOne(ctx, filter).Decode(&context)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &context, nil
}
