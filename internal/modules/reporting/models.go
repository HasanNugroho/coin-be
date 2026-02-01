package reporting

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DailyFinancialReport struct {
	ID                          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID                      primitive.ObjectID `bson:"user_id" json:"user_id"`
	ReportDate                  time.Time          `bson:"report_date" json:"report_date"`
	OpeningBalance              float64            `bson:"opening_balance" json:"opening_balance"`
	ClosingBalance              float64            `bson:"closing_balance" json:"closing_balance"`
	TotalIncome                 float64            `bson:"total_income" json:"total_income"`
	TotalExpense                float64            `bson:"total_expense" json:"total_expense"`
	TotalTransferIn             float64            `bson:"total_transfer_in" json:"total_transfer_in"`
	TotalTransferOut            float64            `bson:"total_transfer_out" json:"total_transfer_out"`
	ExpenseByCategory           []CategoryBreakdown `bson:"expense_by_category" json:"expense_by_category"`
	TransactionsGroupedByPocket []PocketBreakdown   `bson:"transactions_grouped_by_pocket" json:"transactions_grouped_by_pocket"`
	GeneratedAt                 time.Time          `bson:"generated_at" json:"generated_at"`
	IsFinal                     bool               `bson:"is_final" json:"is_final"`
	CreatedAt                   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt                   time.Time          `bson:"updated_at" json:"updated_at"`
}

type CategoryBreakdown struct {
	CategoryID       primitive.ObjectID `bson:"category_id" json:"category_id"`
	CategoryName     string             `bson:"category_name" json:"category_name"`
	Amount           float64            `bson:"amount" json:"amount"`
	TransactionCount int32              `bson:"transaction_count" json:"transaction_count"`
}

type PocketBreakdown struct {
	PocketID         primitive.ObjectID `bson:"pocket_id" json:"pocket_id"`
	PocketName       string             `bson:"pocket_name" json:"pocket_name"`
	PocketType       string             `bson:"pocket_type" json:"pocket_type"`
	PocketBalance    float64            `bson:"pocket_balance" json:"pocket_balance"`
	Income           float64            `bson:"income" json:"income"`
	Expense          float64            `bson:"expense" json:"expense"`
	TransferIn       float64            `bson:"transfer_in" json:"transfer_in"`
	TransferOut      float64            `bson:"transfer_out" json:"transfer_out"`
	TransactionCount int32              `bson:"transaction_count" json:"transaction_count"`
}

type DailyFinancialSnapshot struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID          primitive.ObjectID `bson:"user_id" json:"user_id"`
	SnapshotDate    time.Time          `bson:"snapshot_date" json:"snapshot_date"`
	PocketBalances  []PocketSnapshot   `bson:"pocket_balances" json:"pocket_balances"`
	TotalBalance    float64            `bson:"total_balance" json:"total_balance"`
	FreeMoneyTotal  float64            `bson:"free_money_total" json:"free_money_total"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}

type PocketSnapshot struct {
	PocketID   primitive.ObjectID `bson:"pocket_id" json:"pocket_id"`
	PocketName string             `bson:"pocket_name" json:"pocket_name"`
	PocketType string             `bson:"pocket_type" json:"pocket_type"`
	Balance    float64            `bson:"balance" json:"balance"`
	Currency   string             `bson:"currency" json:"currency"`
}

type MonthlyFinancialSummary struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID             primitive.ObjectID `bson:"user_id" json:"user_id"`
	YearMonth          string             `bson:"year_month" json:"year_month"`
	TotalIncome        float64            `bson:"total_income" json:"total_income"`
	TotalExpense       float64            `bson:"total_expense" json:"total_expense"`
	TotalTransferIn    float64            `bson:"total_transfer_in" json:"total_transfer_in"`
	TotalTransferOut   float64            `bson:"total_transfer_out" json:"total_transfer_out"`
	ExpenseByCategory  []CategoryBreakdown `bson:"expense_by_category" json:"expense_by_category"`
	IncomeByCategory   []CategoryBreakdown `bson:"income_by_category" json:"income_by_category"`
	PocketSummary      []PocketSummary     `bson:"pocket_summary" json:"pocket_summary"`
	TransactionCount   int32              `bson:"transaction_count" json:"transaction_count"`
	GeneratedAt        time.Time          `bson:"generated_at" json:"generated_at"`
	IsFinal            bool               `bson:"is_final" json:"is_final"`
	CreatedAt          time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt          time.Time          `bson:"updated_at" json:"updated_at"`
}

type PocketSummary struct {
	PocketID         primitive.ObjectID `bson:"pocket_id" json:"pocket_id"`
	PocketName       string             `bson:"pocket_name" json:"pocket_name"`
	PocketType       string             `bson:"pocket_type" json:"pocket_type"`
	Income           float64            `bson:"income" json:"income"`
	Expense          float64            `bson:"expense" json:"expense"`
	TransferIn       float64            `bson:"transfer_in" json:"transfer_in"`
	TransferOut      float64            `bson:"transfer_out" json:"transfer_out"`
	TransactionCount int32              `bson:"transaction_count" json:"transaction_count"`
}

type PocketBalanceSnapshot struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID            primitive.ObjectID `bson:"user_id" json:"user_id"`
	PocketID          primitive.ObjectID `bson:"pocket_id" json:"pocket_id"`
	SnapshotDate      time.Time          `bson:"snapshot_date" json:"snapshot_date"`
	Balance           float64            `bson:"balance" json:"balance"`
	BalanceChange     float64            `bson:"balance_change" json:"balance_change"`
	DailyIncome       float64            `bson:"daily_income" json:"daily_income"`
	DailyExpense      float64            `bson:"daily_expense" json:"daily_expense"`
	DailyTransferIn   float64            `bson:"daily_transfer_in" json:"daily_transfer_in"`
	DailyTransferOut  float64            `bson:"daily_transfer_out" json:"daily_transfer_out"`
	CreatedAt         time.Time          `bson:"created_at" json:"created_at"`
}

type AIFinancialContext struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID            primitive.ObjectID `bson:"user_id" json:"user_id"`
	ContextDate       time.Time          `bson:"context_date" json:"context_date"`
	CurrentBalance    float64            `bson:"current_balance" json:"current_balance"`
	FreeMoney         float64            `bson:"free_money" json:"free_money"`
	Last30Days        Last30DaysMetrics  `bson:"last_30_days" json:"last_30_days"`
	YearToDate        YearToDateMetrics  `bson:"year_to_date" json:"year_to_date"`
	Pockets           []AIPocketData     `bson:"pockets" json:"pockets"`
	SpendingPatterns  SpendingPatterns   `bson:"spending_patterns" json:"spending_patterns"`
	Alerts            []Alert            `bson:"alerts" json:"alerts"`
	UpdatedAt         time.Time          `bson:"updated_at" json:"updated_at"`
}

type Last30DaysMetrics struct {
	TotalIncome              float64                `bson:"total_income" json:"total_income"`
	TotalExpense             float64                `bson:"total_expense" json:"total_expense"`
	NetChange                float64                `bson:"net_change" json:"net_change"`
	AverageDailyExpense      float64                `bson:"average_daily_expense" json:"average_daily_expense"`
	TransactionCount         int32                  `bson:"transaction_count" json:"transaction_count"`
	TopExpenseCategories     []TopCategory          `bson:"top_expense_categories" json:"top_expense_categories"`
}

type YearToDateMetrics struct {
	TotalIncome             float64 `bson:"total_income" json:"total_income"`
	TotalExpense            float64 `bson:"total_expense" json:"total_expense"`
	NetChange               float64 `bson:"net_change" json:"net_change"`
	AverageMonthlyExpense   float64 `bson:"average_monthly_expense" json:"average_monthly_expense"`
}

type TopCategory struct {
	CategoryName string  `bson:"category_name" json:"category_name"`
	Amount       float64 `bson:"amount" json:"amount"`
	Percentage   float64 `bson:"percentage" json:"percentage"`
}

type AIPocketData struct {
	PocketID         primitive.ObjectID `bson:"pocket_id" json:"pocket_id"`
	PocketName       string             `bson:"pocket_name" json:"pocket_name"`
	PocketType       string             `bson:"pocket_type" json:"pocket_type"`
	Balance          float64            `bson:"balance" json:"balance"`
	PercentageTotal  float64            `bson:"percentage_of_total" json:"percentage_of_total"`
	MonthlyTrend     []MonthlyBalance   `bson:"monthly_trend" json:"monthly_trend"`
}

type MonthlyBalance struct {
	Month   string  `bson:"month" json:"month"`
	Balance float64 `bson:"balance" json:"balance"`
}

type SpendingPatterns struct {
	HighestExpenseDayOfWeek string  `bson:"highest_expense_day_of_week" json:"highest_expense_day_of_week"`
	HighestExpenseCategory  string  `bson:"highest_expense_category" json:"highest_expense_category"`
	AverageTransactionAmount float64 `bson:"average_transaction_amount" json:"average_transaction_amount"`
	LargestTransaction      float64 `bson:"largest_transaction" json:"largest_transaction"`
	SmallestTransaction     float64 `bson:"smallest_transaction" json:"smallest_transaction"`
}

type Alert struct {
	Type     string `bson:"type" json:"type"`
	Message  string `bson:"message" json:"message"`
	Severity string `bson:"severity" json:"severity"`
}
