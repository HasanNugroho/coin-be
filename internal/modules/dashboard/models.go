package dashboard

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DailySummary struct {
	ID                primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	UserID            primitive.ObjectID  `bson:"user_id" json:"user_id"`
	Date              time.Time           `bson:"date" json:"date"`
	TotalIncome       float64             `bson:"total_income" json:"total_income"`
	TotalExpense      float64             `bson:"total_expense" json:"total_expense"`
	CategoryBreakdown []CategoryBreakdown `bson:"category_breakdown" json:"category_breakdown"`
	CreatedAt         time.Time           `bson:"created_at" json:"created_at"`
}

type CategoryBreakdown struct {
	CategoryID   *primitive.ObjectID `bson:"category_id,omitempty" json:"category_id,omitempty"`
	CategoryName string              `bson:"category_name" json:"category_name"`
	Type         string              `bson:"type" json:"type"`
	Amount       float64             `bson:"amount" json:"amount"`
}

type DashboardSummary struct {
	TotalNetWorth float64   `json:"total_net_worth"`
	PeriodIncome  float64   `json:"period_income"`
	PeriodExpense float64   `json:"period_expense"`
	PeriodNet     float64   `json:"period_net"`
	TimeRange     TimeRange `json:"time_range"`
}

type ChartDataPoint struct {
	Date    string  `json:"date"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
}

type CategoryChartData struct {
	CategoryID   string  `json:"category_id,omitempty"`
	CategoryName string  `json:"category_name"`
	Amount       float64 `json:"amount"`
	Percentage   float64 `json:"percentage"`
}

type DashboardCharts struct {
	CashFlowTrend    []ChartDataPoint    `json:"cash_flow_trend"`
	IncomeBreakdown  []CategoryChartData `json:"income_breakdown"`
	ExpenseBreakdown []CategoryChartData `json:"expense_breakdown"`
}
