package report

import "go.mongodb.org/mongo-driver/bson/primitive"

type DashboardSummary struct {
	TotalBalance       float64 `json:"total_balance"`
	FreeCash           float64 `json:"free_cash"`
	IncomeThisMonth    float64 `json:"income_this_month"`
	ExpenseThisMonth   float64 `json:"expense_this_month"`
	RemainingThisMonth float64 `json:"remaining_this_month"`
}

type IncomeReport struct {
	TotalIncome float64             `json:"total_income"`
	ByCategory  []CategoryBreakdown `json:"by_category"`
	ByMonth     []MonthlyBreakdown  `json:"by_month"`
}

type ExpenseReport struct {
	TotalExpense float64               `json:"total_expense"`
	ByCategory   []CategoryBreakdown   `json:"by_category"`
	ByAllocation []AllocationBreakdown `json:"by_allocation"`
	ByMonth      []MonthlyBreakdown    `json:"by_month"`
}

type AllocationReport struct {
	AllocationID      primitive.ObjectID `json:"allocation_id"`
	AllocationName    string             `json:"allocation_name"`
	CurrentBalance    float64            `json:"current_balance"`
	TargetAmount      *float64           `json:"target_amount,omitempty"`
	Progress          float64            `json:"progress"`
	TotalAllocated    float64            `json:"total_allocated"`
	TotalSpent        float64            `json:"total_spent"`
	DistributionCount int                `json:"distribution_count"`
}

type TargetProgress struct {
	TargetID      primitive.ObjectID `json:"target_id"`
	TargetName    string             `json:"target_name"`
	TargetAmount  float64            `json:"target_amount"`
	CurrentAmount float64            `json:"current_amount"`
	Progress      float64            `json:"progress"`
	Status        string             `json:"status"`
	DaysRemaining int                `json:"days_remaining"`
}

type CategoryBreakdown struct {
	CategoryID   primitive.ObjectID `json:"category_id"`
	CategoryName string             `json:"category_name"`
	Amount       float64            `json:"amount"`
	Percentage   float64            `json:"percentage"`
	Count        int                `json:"count"`
}

type AllocationBreakdown struct {
	AllocationID   primitive.ObjectID `json:"allocation_id"`
	AllocationName string             `json:"allocation_name"`
	Amount         float64            `json:"amount"`
	Percentage     float64            `json:"percentage"`
	Count          int                `json:"count"`
}

type MonthlyBreakdown struct {
	Month  string  `json:"month"`
	Year   int     `json:"year"`
	Amount float64 `json:"amount"`
	Count  int     `json:"count"`
}
