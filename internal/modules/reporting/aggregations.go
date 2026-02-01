package reporting

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AggregationHelper provides pre-built aggregation pipelines for dashboard and reporting
type AggregationHelper struct {
	db *mongo.Database
}

func NewAggregationHelper(db *mongo.Database) *AggregationHelper {
	return &AggregationHelper{db: db}
}

// ============================================================================
// KPI AGGREGATIONS
// ============================================================================

// GetTotalBalance aggregates balance across all active pockets
func (a *AggregationHelper) GetTotalBalance(ctx context.Context, userID primitive.ObjectID) (primitive.Decimal128, error) {
	pipeline := mongo.Pipeline{
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "user_id", Value: userID},
				{Key: "is_active", Value: true},
				{Key: "deleted_at", Value: nil},
			}},
		},
		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: nil},
				{Key: "total_balance", Value: bson.D{
					{Key: "$sum", Value: "$balance"},
				}},
			}},
		},
		bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "total_balance", Value: 1},
			}},
		},
	}

	cursor, err := a.db.Collection("pockets").Aggregate(ctx, pipeline)
	if err != nil {
		return primitive.Decimal128{}, err
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		return primitive.Decimal128{}, err
	}

	if len(result) == 0 {
		return primitive.NewDecimal128(0, 0), nil
	}

	return result[0]["total_balance"].(primitive.Decimal128), nil
}

// GetMonthlyIncome aggregates income for a specific month from daily reports
func (a *AggregationHelper) GetMonthlyIncome(ctx context.Context, userID primitive.ObjectID, month time.Time) (primitive.Decimal128, error) {
	startOfMonth := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	pipeline := mongo.Pipeline{
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "user_id", Value: userID},
				{Key: "report_date", Value: bson.D{
					{Key: "$gte", Value: startOfMonth},
					{Key: "$lte", Value: endOfMonth},
				}},
				{Key: "is_final", Value: true},
			}},
		},
		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: nil},
				{Key: "total_income", Value: bson.D{
					{Key: "$sum", Value: "$total_income"},
				}},
			}},
		},
	}

	cursor, err := a.db.Collection("daily_financial_reports").Aggregate(ctx, pipeline)
	if err != nil {
		return primitive.Decimal128{}, err
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		return primitive.Decimal128{}, err
	}

	if len(result) == 0 {
		return primitive.NewDecimal128(0, 0), nil
	}

	return result[0]["total_income"].(primitive.Decimal128), nil
}

// GetMonthlyExpense aggregates expense for a specific month from daily reports
func (a *AggregationHelper) GetMonthlyExpense(ctx context.Context, userID primitive.ObjectID, month time.Time) (primitive.Decimal128, error) {
	startOfMonth := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	pipeline := mongo.Pipeline{
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "user_id", Value: userID},
				{Key: "report_date", Value: bson.D{
					{Key: "$gte", Value: startOfMonth},
					{Key: "$lte", Value: endOfMonth},
				}},
				{Key: "is_final", Value: true},
			}},
		},
		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: nil},
				{Key: "total_expense", Value: bson.D{
					{Key: "$sum", Value: "$total_expense"},
				}},
			}},
		},
	}

	cursor, err := a.db.Collection("daily_financial_reports").Aggregate(ctx, pipeline)
	if err != nil {
		return primitive.Decimal128{}, err
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		return primitive.Decimal128{}, err
	}

	if len(result) == 0 {
		return primitive.NewDecimal128(0, 0), nil
	}

	return result[0]["total_expense"].(primitive.Decimal128), nil
}

// GetFreeMoneyTotal aggregates balance for main and allocation pockets
func (a *AggregationHelper) GetFreeMoneyTotal(ctx context.Context, userID primitive.ObjectID) (primitive.Decimal128, error) {
	pipeline := mongo.Pipeline{
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "user_id", Value: userID},
				{Key: "type", Value: bson.D{
					{Key: "$in", Value: []string{"main", "allocation"}},
				}},
				{Key: "is_active", Value: true},
				{Key: "deleted_at", Value: nil},
			}},
		},
		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: nil},
				{Key: "free_money", Value: bson.D{
					{Key: "$sum", Value: "$balance"},
				}},
			}},
		},
	}

	cursor, err := a.db.Collection("pockets").Aggregate(ctx, pipeline)
	if err != nil {
		return primitive.Decimal128{}, err
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		return primitive.Decimal128{}, err
	}

	if len(result) == 0 {
		return primitive.NewDecimal128(0, 0), nil
	}

	return result[0]["free_money"].(primitive.Decimal128), nil
}

// ============================================================================
// CHART AGGREGATIONS
// ============================================================================

// GetMonthlyIncomeExpenseChart returns 12 months of income vs expense
func (a *AggregationHelper) GetMonthlyIncomeExpenseChart(ctx context.Context, userID primitive.ObjectID, months int) ([]bson.M, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, -months, 0)

	pipeline := mongo.Pipeline{
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "user_id", Value: userID},
				{Key: "month", Value: bson.D{
					{Key: "$gte", Value: startDate},
					{Key: "$lte", Value: endDate},
				}},
			}},
		},
		bson.D{
			{Key: "$sort", Value: bson.D{
				{Key: "month", Value: 1},
			}},
		},
		bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "month", Value: bson.D{
					{Key: "$dateToString", Value: bson.D{
						{Key: "format", Value: "%Y-%m"},
						{Key: "date", Value: "$month"},
					}},
				}},
				{Key: "income", Value: 1},
				{Key: "expense", Value: 1},
			}},
		},
	}

	cursor, err := a.db.Collection("monthly_financial_summaries").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []bson.M
	return result, cursor.All(ctx, &result)
}

// GetPocketBalanceDistribution returns balance breakdown by pocket type
func (a *AggregationHelper) GetPocketBalanceDistribution(ctx context.Context, userID primitive.ObjectID) ([]bson.M, error) {
	pipeline := mongo.Pipeline{
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "user_id", Value: userID},
				{Key: "is_active", Value: true},
				{Key: "deleted_at", Value: nil},
			}},
		},
		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$type"},
				{Key: "pockets", Value: bson.D{
					{Key: "$push", Value: bson.D{
						{Key: "pocket_id", Value: "$_id"},
						{Key: "pocket_name", Value: "$name"},
						{Key: "balance", Value: "$balance"},
					}},
				}},
				{Key: "total_by_type", Value: bson.D{
					{Key: "$sum", Value: "$balance"},
				}},
			}},
		},
		bson.D{
			{Key: "$sort", Value: bson.D{
				{Key: "total_by_type", Value: -1},
			}},
		},
	}

	cursor, err := a.db.Collection("pockets").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []bson.M
	return result, cursor.All(ctx, &result)
}

// GetExpenseCategoryDistribution returns expense breakdown by category for current month
func (a *AggregationHelper) GetExpenseCategoryDistribution(ctx context.Context, userID primitive.ObjectID, month time.Time) ([]bson.M, error) {
	startOfMonth := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	pipeline := mongo.Pipeline{
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "user_id", Value: userID},
				{Key: "report_date", Value: bson.D{
					{Key: "$gte", Value: startOfMonth},
					{Key: "$lte", Value: endOfMonth},
				}},
				{Key: "is_final", Value: true},
			}},
		},
		bson.D{
			{Key: "$unwind", Value: "$expense_by_category"},
		},
		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$expense_by_category.category_id"},
				{Key: "category_name", Value: bson.D{
					{Key: "$first", Value: "$expense_by_category.category_name"},
				}},
				{Key: "total_amount", Value: bson.D{
					{Key: "$sum", Value: "$expense_by_category.amount"},
				}},
				{Key: "transaction_count", Value: bson.D{
					{Key: "$sum", Value: "$expense_by_category.transaction_count"},
				}},
			}},
		},
		bson.D{
			{Key: "$sort", Value: bson.D{
				{Key: "total_amount", Value: -1},
			}},
		},
	}

	cursor, err := a.db.Collection("daily_financial_reports").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []bson.M
	return result, cursor.All(ctx, &result)
}

// ============================================================================
// AI AGGREGATIONS
// ============================================================================

// GetAnomalySummary returns summary of anomalies for the user
func (a *AggregationHelper) GetAnomalySummary(ctx context.Context, userID primitive.ObjectID, days int) (bson.M, error) {
	startDate := time.Now().AddDate(0, 0, -days)

	pipeline := mongo.Pipeline{
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "user_id", Value: userID},
				{Key: "is_anomaly", Value: true},
				{Key: "date", Value: bson.D{
					{Key: "$gte", Value: startDate},
				}},
				{Key: "deleted_at", Value: nil},
			}},
		},
		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: nil},
				{Key: "count", Value: bson.D{
					{Key: "$sum", Value: 1},
				}},
				{Key: "avg_anomaly_score", Value: bson.D{
					{Key: "$avg", Value: "$anomaly_score"},
				}},
				{Key: "total_anomalous_amount", Value: bson.D{
					{Key: "$sum", Value: "$amount"},
				}},
			}},
		},
	}

	cursor, err := a.db.Collection("ai_transaction_enrichment").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return bson.M{
			"count":                  0,
			"avg_anomaly_score":      0,
			"total_anomalous_amount": 0,
		}, nil
	}

	return result[0], nil
}

// GetSpendingTrendsByCategory returns spending trends for top categories
func (a *AggregationHelper) GetSpendingTrendsByCategory(ctx context.Context, userID primitive.ObjectID, limit int) ([]bson.M, error) {
	pipeline := mongo.Pipeline{
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "user_id", Value: userID},
				{Key: "deleted_at", Value: nil},
			}},
		},
		bson.D{
			{Key: "$sort", Value: bson.D{
				{Key: "frequency_per_month", Value: -1},
			}},
		},
		bson.D{
			{Key: "$limit", Value: int64(limit)},
		},
		bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "category", Value: 1},
				{Key: "merchant", Value: 1},
				{Key: "avg_amount", Value: 1},
				{Key: "frequency_per_month", Value: 1},
				{Key: "trend", Value: 1},
				{Key: "trend_percent_change", Value: 1},
				{Key: "is_seasonal", Value: 1},
				{Key: "is_essential", Value: 1},
			}},
		},
	}

	cursor, err := a.db.Collection("ai_spending_patterns").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []bson.M
	return result, cursor.All(ctx, &result)
}

// GetRecurringExpensesSummary returns summary of recurring expenses
func (a *AggregationHelper) GetRecurringExpensesSummary(ctx context.Context, userID primitive.ObjectID) (primitive.Decimal128, error) {
	pipeline := mongo.Pipeline{
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "user_id", Value: userID},
				{Key: "is_recurring", Value: true},
				{Key: "deleted_at", Value: nil},
			}},
		},
		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: nil},
				{Key: "total_recurring", Value: bson.D{
					{Key: "$sum", Value: bson.D{
						{Key: "$multiply", Value: []interface{}{
							"$avg_amount",
							"$frequency_per_month",
						}},
					}},
				}},
			}},
		},
	}

	cursor, err := a.db.Collection("ai_spending_patterns").Aggregate(ctx, pipeline)
	if err != nil {
		return primitive.Decimal128{}, err
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		return primitive.Decimal128{}, err
	}

	if len(result) == 0 {
		return primitive.NewDecimal128(0, 0), nil
	}

	return result[0]["total_recurring"].(primitive.Decimal128), nil
}
