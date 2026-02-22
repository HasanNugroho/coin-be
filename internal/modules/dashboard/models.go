package dashboard

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
