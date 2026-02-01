package reporting

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DailyFinancialReport represents a pre-aggregated daily financial report
type DailyFinancialReport struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	ReportDate time.Time          `bson:"report_date" json:"report_date"`

	// Opening/Closing Balances
	OpeningBalance primitive.Decimal128 `bson:"opening_balance" json:"opening_balance"`
	ClosingBalance primitive.Decimal128 `bson:"closing_balance" json:"closing_balance"`

	// Daily Aggregates
	TotalIncome      primitive.Decimal128 `bson:"total_income" json:"total_income"`
	TotalExpense     primitive.Decimal128 `bson:"total_expense" json:"total_expense"`
	TotalTransferIn  primitive.Decimal128 `bson:"total_transfer_in" json:"total_transfer_in"`
	TotalTransferOut primitive.Decimal128 `bson:"total_transfer_out" json:"total_transfer_out"`

	// Expense Breakdown by Category
	ExpenseByCategory []ExpenseByCategory `bson:"expense_by_category" json:"expense_by_category"`

	// Transactions Grouped by Pocket
	TransactionsByPocket []TransactionsByPocket `bson:"transactions_by_pocket" json:"transactions_by_pocket"`

	// Metadata
	GeneratedAt time.Time `bson:"generated_at" json:"generated_at"`
	IsFinal     bool      `bson:"is_final" json:"is_final"`

	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

type ExpenseByCategory struct {
	CategoryID       primitive.ObjectID   `bson:"category_id" json:"category_id"`
	CategoryName     string               `bson:"category_name" json:"category_name"`
	Amount           primitive.Decimal128 `bson:"amount" json:"amount"`
	TransactionCount int32                `bson:"transaction_count" json:"transaction_count"`
}

type TransactionsByPocket struct {
	PocketID           primitive.ObjectID   `bson:"pocket_id" json:"pocket_id"`
	PocketName         string               `bson:"pocket_name" json:"pocket_name"`
	PocketType         string               `bson:"pocket_type" json:"pocket_type"`
	Income             primitive.Decimal128 `bson:"income" json:"income"`
	Expense            primitive.Decimal128 `bson:"expense" json:"expense"`
	TransferIn         primitive.Decimal128 `bson:"transfer_in" json:"transfer_in"`
	TransferOut        primitive.Decimal128 `bson:"transfer_out" json:"transfer_out"`
	OpeningBalance     primitive.Decimal128 `bson:"opening_balance" json:"opening_balance"`
	ClosingBalance     primitive.Decimal128 `bson:"closing_balance" json:"closing_balance"`
	TransactionCount   int32                `bson:"transaction_count" json:"transaction_count"`
	Transactions       []ReportTransaction  `bson:"transactions" json:"transactions"`
}

type ReportTransaction struct {
	TransactionID primitive.ObjectID   `bson:"transaction_id" json:"transaction_id"`
	Type          string               `bson:"type" json:"type"`
	Amount        primitive.Decimal128 `bson:"amount" json:"amount"`
	CategoryID    *primitive.ObjectID  `bson:"category_id,omitempty" json:"category_id,omitempty"`
	CategoryName  string               `bson:"category_name" json:"category_name"`
	Note          *string              `bson:"note,omitempty" json:"note,omitempty"`
	Timestamp     time.Time            `bson:"timestamp" json:"timestamp"`
}

// DailyFinancialSnapshot represents a point-in-time snapshot of daily financial state
type DailyFinancialSnapshot struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	UserID       primitive.ObjectID `bson:"user_id" json:"user_id"`
	SnapshotDate time.Time          `bson:"snapshot_date" json:"snapshot_date"`

	// Pocket balances at end of day
	PocketBalances []PocketBalanceSnapshot `bson:"pocket_balances" json:"pocket_balances"`

	// Daily totals
	TotalBalance     primitive.Decimal128 `bson:"total_balance" json:"total_balance"`
	TotalIncome      primitive.Decimal128 `bson:"total_income" json:"total_income"`
	TotalExpense     primitive.Decimal128 `bson:"total_expense" json:"total_expense"`
	TotalTransferIn  primitive.Decimal128 `bson:"total_transfer_in" json:"total_transfer_in"`
	TotalTransferOut primitive.Decimal128 `bson:"total_transfer_out" json:"total_transfer_out"`

	// Cumulative (year-to-date)
	YTDIncome  primitive.Decimal128 `bson:"ytd_income" json:"ytd_income"`
	YTDExpense primitive.Decimal128 `bson:"ytd_expense" json:"ytd_expense"`
	YTDNet     primitive.Decimal128 `bson:"ytd_net" json:"ytd_net"`

	// Metadata
	TransactionCount int32     `bson:"transaction_count" json:"transaction_count"`
	IsComplete       bool      `bson:"is_complete" json:"is_complete"`
	GeneratedAt      time.Time `bson:"generated_at" json:"generated_at"`

	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

type PocketBalanceSnapshot struct {
	PocketID      primitive.ObjectID   `bson:"pocket_id" json:"pocket_id"`
	PocketName    string               `bson:"pocket_name" json:"pocket_name"`
	PocketType    string               `bson:"pocket_type" json:"pocket_type"`
	Balance       primitive.Decimal128 `bson:"balance" json:"balance"`
	BalanceChange primitive.Decimal128 `bson:"balance_change" json:"balance_change"`
}

// MonthlyFinancialSummary represents pre-aggregated monthly financial data
type MonthlyFinancialSummary struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
	Month  time.Time          `bson:"month" json:"month"`

	// Monthly totals
	Income       primitive.Decimal128 `bson:"income" json:"income"`
	Expense      primitive.Decimal128 `bson:"expense" json:"expense"`
	TransferIn   primitive.Decimal128 `bson:"transfer_in" json:"transfer_in"`
	TransferOut  primitive.Decimal128 `bson:"transfer_out" json:"transfer_out"`
	Net          primitive.Decimal128 `bson:"net" json:"net"`

	// Opening/Closing
	OpeningBalance primitive.Decimal128 `bson:"opening_balance" json:"opening_balance"`
	ClosingBalance primitive.Decimal128 `bson:"closing_balance" json:"closing_balance"`

	// By category
	ExpenseByCategory []MonthlyCategoryBreakdown `bson:"expense_by_category" json:"expense_by_category"`

	// By pocket
	ByPocket []MonthlyPocketBreakdown `bson:"by_pocket" json:"by_pocket"`

	// Cumulative
	YTDIncome  primitive.Decimal128 `bson:"ytd_income" json:"ytd_income"`
	YTDExpense primitive.Decimal128 `bson:"ytd_expense" json:"ytd_expense"`
	YTDNet     primitive.Decimal128 `bson:"ytd_net" json:"ytd_net"`

	// Metadata
	DayCount         int32     `bson:"day_count" json:"day_count"`
	TransactionCount int32     `bson:"transaction_count" json:"transaction_count"`
	IsComplete       bool      `bson:"is_complete" json:"is_complete"`
	CreatedAt        time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt        *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

type MonthlyCategoryBreakdown struct {
	CategoryID       primitive.ObjectID   `bson:"category_id" json:"category_id"`
	CategoryName     string               `bson:"category_name" json:"category_name"`
	Amount           primitive.Decimal128 `bson:"amount" json:"amount"`
	PercentOfTotal   float64              `bson:"percent_of_total" json:"percent_of_total"`
	TransactionCount int32                `bson:"transaction_count" json:"transaction_count"`
}

type MonthlyPocketBreakdown struct {
	PocketID   primitive.ObjectID   `bson:"pocket_id" json:"pocket_id"`
	PocketName string               `bson:"pocket_name" json:"pocket_name"`
	PocketType string               `bson:"pocket_type" json:"pocket_type"`
	Income     primitive.Decimal128 `bson:"income" json:"income"`
	Expense    primitive.Decimal128 `bson:"expense" json:"expense"`
	Net        primitive.Decimal128 `bson:"net" json:"net"`
}

// PocketBalanceSnapshot represents a point-in-time balance record for audit trail
type PocketBalanceHistorySnapshot struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	UserID   primitive.ObjectID `bson:"user_id" json:"user_id"`
	PocketID primitive.ObjectID `bson:"pocket_id" json:"pocket_id"`

	// Balance at snapshot time
	Balance       primitive.Decimal128 `bson:"balance" json:"balance"`
	BalanceBefore primitive.Decimal128 `bson:"balance_before" json:"balance_before"`
	Change        primitive.Decimal128 `bson:"change" json:"change"`

	// Snapshot metadata
	SnapshotTime time.Time `bson:"snapshot_time" json:"snapshot_time"`
	SnapshotType string    `bson:"snapshot_type" json:"snapshot_type"` // "hourly", "daily", "transaction"

	// Transaction that caused change (if applicable)
	TransactionID     *primitive.ObjectID `bson:"transaction_id,omitempty" json:"transaction_id,omitempty"`
	TransactionType   *string             `bson:"transaction_type,omitempty" json:"transaction_type,omitempty"`
	TransactionAmount *primitive.Decimal128 `bson:"transaction_amount,omitempty" json:"transaction_amount,omitempty"`

	// Context
	PocketName string `bson:"pocket_name" json:"pocket_name"`
	PocketType string `bson:"pocket_type" json:"pocket_type"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
