package reporting

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AITransactionEnrichment represents enriched transaction data for AI/chatbot consumption
type AITransactionEnrichment struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
	TransactionID primitive.ObjectID `bson:"transaction_id" json:"transaction_id"`

	// Original transaction data (denormalized for LLM)
	TransactionType string               `bson:"transaction_type" json:"transaction_type"`
	Amount          primitive.Decimal128 `bson:"amount" json:"amount"`
	Date            time.Time            `bson:"date" json:"date"`

	// Enrichment
	MerchantName      string  `bson:"merchant_name" json:"merchant_name"`
	MerchantCategory  string  `bson:"merchant_category" json:"merchant_category"`
	ConfidenceScore   float64 `bson:"confidence_score" json:"confidence_score"` // 0-1

	// NLP-friendly description
	Description string   `bson:"description" json:"description"`
	Tags        []string `bson:"tags" json:"tags"`

	// Anomaly detection
	IsAnomaly      bool    `bson:"is_anomaly" json:"is_anomaly"`
	AnomalyReason  string  `bson:"anomaly_reason" json:"anomaly_reason"`
	AnomalyScore   float64 `bson:"anomaly_score" json:"anomaly_score"` // 0-1

	// Spending pattern context
	CategoryAvgAmount     primitive.Decimal128 `bson:"category_avg_amount" json:"category_avg_amount"`
	CategoryAvgFrequency  int32                `bson:"category_avg_frequency" json:"category_avg_frequency"`
	IsRecurring           bool                 `bson:"is_recurring" json:"is_recurring"`

	// Budget context
	BudgetCategory         string               `bson:"budget_category" json:"budget_category"`
	BudgetRemaining        primitive.Decimal128 `bson:"budget_remaining" json:"budget_remaining"`
	BudgetUtilizationPct   float64              `bson:"budget_utilization_percent" json:"budget_utilization_percent"`

	// Temporal context
	DayOfWeek   string `bson:"day_of_week" json:"day_of_week"`
	IsWeekend   bool   `bson:"is_weekend" json:"is_weekend"`
	IsHoliday   bool   `bson:"is_holiday" json:"is_holiday"`

	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

// AISpendingPattern represents identified spending patterns for AI analysis
type AISpendingPattern struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	UserID   primitive.ObjectID `bson:"user_id" json:"user_id"`
	Category string             `bson:"category" json:"category"`
	Merchant string             `bson:"merchant" json:"merchant"`

	// Statistics
	AvgAmount    primitive.Decimal128 `bson:"avg_amount" json:"avg_amount"`
	MedianAmount primitive.Decimal128 `bson:"median_amount" json:"median_amount"`
	StdDev       primitive.Decimal128 `bson:"std_dev" json:"std_dev"`
	MinAmount    primitive.Decimal128 `bson:"min_amount" json:"min_amount"`
	MaxAmount    primitive.Decimal128 `bson:"max_amount" json:"max_amount"`

	// Frequency
	FrequencyPerMonth float64   `bson:"frequency_per_month" json:"frequency_per_month"`
	FrequencyPerWeek  float64   `bson:"frequency_per_week" json:"frequency_per_week"`
	LastTransactionDate time.Time `bson:"last_transaction_date" json:"last_transaction_date"`

	// Trend
	Trend               string  `bson:"trend" json:"trend"` // "increasing", "decreasing", "stable"
	TrendPercentChange  float64 `bson:"trend_percent_change" json:"trend_percent_change"`

	// Seasonality
	IsSeasonal      bool   `bson:"is_seasonal" json:"is_seasonal"`
	SeasonalMonths  []int32 `bson:"seasonal_months" json:"seasonal_months"` // [1, 12] for Jan, Dec

	// Behavioral insights
	PreferredDayOfWeek string `bson:"preferred_day_of_week" json:"preferred_day_of_week"`
	PreferredTimeOfDay string `bson:"preferred_time_of_day" json:"preferred_time_of_day"`

	// AI context
	SpendingCategoryRank int32   `bson:"spending_category_rank" json:"spending_category_rank"`
	IsEssential          bool    `bson:"is_essential" json:"is_essential"`

	DataPoints int32   `bson:"data_points" json:"data_points"`
	Confidence float64 `bson:"confidence" json:"confidence"` // 0-1

	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

// AIFinancialInsight represents AI-generated insights for chatbot
type AIFinancialInsight struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`

	// Insight metadata
	InsightType string `bson:"insight_type" json:"insight_type"` // "anomaly", "trend", "opportunity", "warning"
	Category    string `bson:"category" json:"category"`         // "spending", "income", "savings", "budget"
	Severity    string `bson:"severity" json:"severity"`         // "info", "warning", "critical"

	// Content for chatbot
	Title          string `bson:"title" json:"title"`
	Description    string `bson:"description" json:"description"`
	Recommendation string `bson:"recommendation" json:"recommendation"`

	// Data backing the insight
	MetricName       string               `bson:"metric_name" json:"metric_name"`
	MetricValue      primitive.Decimal128 `bson:"metric_value" json:"metric_value"`
	MetricBaseline   primitive.Decimal128 `bson:"metric_baseline" json:"metric_baseline"`
	MetricChangePct  float64              `bson:"metric_change_percent" json:"metric_change_percent"`

	// Context
	AffectedCategories []string             `bson:"affected_categories" json:"affected_categories"`
	AffectedPockets    []primitive.ObjectID `bson:"affected_pockets" json:"affected_pockets"`
	DateRange          DateRange            `bson:"date_range" json:"date_range"`

	// Engagement
	IsActionable bool   `bson:"is_actionable" json:"is_actionable"`
	ActionURL    string `bson:"action_url" json:"action_url"`

	// Lifecycle
	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	ExpiresAt time.Time  `bson:"expires_at" json:"expires_at"`
	IsRead    bool       `bson:"is_read" json:"is_read"`
	ReadAt    *time.Time `bson:"read_at,omitempty" json:"read_at,omitempty"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

type DateRange struct {
	Start time.Time `bson:"start" json:"start"`
	End   time.Time `bson:"end" json:"end"`
}

// AIFinancialContext represents user's financial context for LLM consumption
type AIFinancialContext struct {
	UserID   primitive.ObjectID `bson:"user_id" json:"user_id"`
	Period   string             `bson:"period" json:"period"` // "last_30_days", "last_month", etc.
	Summary  FinancialSummary   `bson:"summary" json:"summary"`
	TopSpending []SpendingCategory `bson:"top_spending_categories" json:"top_spending_categories"`
	Anomalies []AnomalyInfo      `bson:"anomalies" json:"anomalies"`
	BudgetStatus map[string]BudgetStatus `bson:"budget_status" json:"budget_status"`
	Insights []InsightInfo      `bson:"insights" json:"insights"`
	GeneratedAt time.Time        `bson:"generated_at" json:"generated_at"`
}

type FinancialSummary struct {
	TotalIncome     primitive.Decimal128 `bson:"total_income" json:"total_income"`
	TotalExpense    primitive.Decimal128 `bson:"total_expense" json:"total_expense"`
	NetSavings      primitive.Decimal128 `bson:"net_savings" json:"net_savings"`
	SavingsRate     float64              `bson:"savings_rate" json:"savings_rate"`
}

type SpendingCategory struct {
	Name            string               `bson:"name" json:"name"`
	Amount          primitive.Decimal128 `bson:"amount" json:"amount"`
	PercentOfTotal  float64              `bson:"percent_of_total" json:"percent_of_total"`
	Trend           string               `bson:"trend" json:"trend"`
	VsAverage       string               `bson:"vs_average" json:"vs_average"`
}

type AnomalyInfo struct {
	Date        time.Time            `bson:"date" json:"date"`
	Type        string               `bson:"type" json:"type"`
	Description string               `bson:"description" json:"description"`
	Category    string               `bson:"category" json:"category"`
	Severity    string               `bson:"severity" json:"severity"`
}

type BudgetStatus struct {
	Allocated primitive.Decimal128 `bson:"allocated" json:"allocated"`
	Spent     primitive.Decimal128 `bson:"spent" json:"spent"`
	Remaining primitive.Decimal128 `bson:"remaining" json:"remaining"`
}

type InsightInfo struct {
	Type              string               `bson:"type" json:"type"`
	Message           string               `bson:"message" json:"message"`
	PotentialSavings primitive.Decimal128 `bson:"potential_savings" json:"potential_savings"`
}
