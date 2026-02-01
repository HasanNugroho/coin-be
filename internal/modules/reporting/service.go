package reporting

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo              *Repository
	aggregationHelper *AggregationHelper
}

func NewService(repo *Repository, aggregationHelper *AggregationHelper) *Service {
	return &Service{
		repo:              repo,
		aggregationHelper: aggregationHelper,
	}
}

// ============================================================================
// DASHBOARD SERVICE
// ============================================================================

type DashboardKPIs struct {
	TotalBalance     primitive.Decimal128 `json:"total_balance"`
	MonthlyIncome    primitive.Decimal128 `json:"monthly_income"`
	MonthlyExpense   primitive.Decimal128 `json:"monthly_expense"`
	FreeMoneyTotal   primitive.Decimal128 `json:"free_money_total"`
	MonthlyNetChange primitive.Decimal128 `json:"monthly_net_change"`
}

func (s *Service) GetDashboardKPIs(ctx context.Context, userID primitive.ObjectID) (*DashboardKPIs, error) {
	now := time.Now()
	month := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	totalBalance, err := s.aggregationHelper.GetTotalBalance(ctx, userID)
	if err != nil {
		return nil, err
	}

	monthlyIncome, err := s.aggregationHelper.GetMonthlyIncome(ctx, userID, month)
	if err != nil {
		return nil, err
	}

	monthlyExpense, err := s.aggregationHelper.GetMonthlyExpense(ctx, userID, month)
	if err != nil {
		return nil, err
	}

	freeMoneyTotal, err := s.aggregationHelper.GetFreeMoneyTotal(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &DashboardKPIs{
		TotalBalance:     totalBalance,
		MonthlyIncome:    monthlyIncome,
		MonthlyExpense:   monthlyExpense,
		FreeMoneyTotal:   freeMoneyTotal,
		MonthlyNetChange: subtractDecimal128(monthlyIncome, monthlyExpense),
	}, nil
}

type DashboardCharts struct {
	IncomeExpenseChart   interface{} `json:"income_expense_chart"`
	PocketDistribution   interface{} `json:"pocket_distribution"`
	CategoryDistribution interface{} `json:"category_distribution"`
}

func (s *Service) GetDashboardCharts(ctx context.Context, userID primitive.ObjectID) (*DashboardCharts, error) {
	incomeExpenseChart, err := s.aggregationHelper.GetMonthlyIncomeExpenseChart(ctx, userID, 12)
	if err != nil {
		return nil, err
	}

	pocketDistribution, err := s.aggregationHelper.GetPocketBalanceDistribution(ctx, userID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	month := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	categoryDistribution, err := s.aggregationHelper.GetExpenseCategoryDistribution(ctx, userID, month)
	if err != nil {
		return nil, err
	}

	return &DashboardCharts{
		IncomeExpenseChart:   incomeExpenseChart,
		PocketDistribution:   pocketDistribution,
		CategoryDistribution: categoryDistribution,
	}, nil
}

// ============================================================================
// REPORTING SERVICE
// ============================================================================

type ReportSummary struct {
	ReportDate       time.Time            `json:"report_date"`
	OpeningBalance   primitive.Decimal128 `json:"opening_balance"`
	ClosingBalance   primitive.Decimal128 `json:"closing_balance"`
	TotalIncome      primitive.Decimal128 `json:"total_income"`
	TotalExpense     primitive.Decimal128 `json:"total_expense"`
	NetChange        primitive.Decimal128 `json:"net_change"`
	TransactionCount int32                `json:"transaction_count"`
}

func (s *Service) GetDailyReportSummary(ctx context.Context, userID primitive.ObjectID, reportDate time.Time) (*ReportSummary, error) {
	report, err := s.repo.GetDailyReport(ctx, userID, reportDate)
	if err != nil {
		return nil, err
	}

	if report == nil {
		return nil, nil
	}

	netChange := subtractDecimal128(report.ClosingBalance, report.OpeningBalance)

	return &ReportSummary{
		ReportDate:       report.ReportDate,
		OpeningBalance:   report.OpeningBalance,
		ClosingBalance:   report.ClosingBalance,
		TotalIncome:      report.TotalIncome,
		TotalExpense:     report.TotalExpense,
		NetChange:        netChange,
		TransactionCount: 0, // Calculated from transactions_by_pocket
	}, nil
}

func (s *Service) GetMonthlyReportSummary(ctx context.Context, userID primitive.ObjectID, month time.Time) (*ReportSummary, error) {
	summary, err := s.repo.GetMonthlySummary(ctx, userID, month)
	if err != nil {
		return nil, err
	}

	if summary == nil {
		return nil, nil
	}

	return &ReportSummary{
		ReportDate:       month,
		OpeningBalance:   summary.OpeningBalance,
		ClosingBalance:   summary.ClosingBalance,
		TotalIncome:      summary.Income,
		TotalExpense:     summary.Expense,
		NetChange:        summary.Net,
		TransactionCount: summary.TransactionCount,
	}, nil
}

// ============================================================================
// AI SERVICE
// ============================================================================

type AnomalySummary struct {
	Count                int                  `json:"count"`
	AvgAnomalyScore      float64              `json:"avg_anomaly_score"`
	TotalAnomalousAmount primitive.Decimal128 `json:"total_anomalous_amount"`
}

func (s *Service) GetAnomalySummary(ctx context.Context, userID primitive.ObjectID, days int) (*AnomalySummary, error) {
	result, err := s.aggregationHelper.GetAnomalySummary(ctx, userID, days)
	if err != nil {
		return nil, err
	}

	count := int(result["count"].(int32))
	avgScore := result["avg_anomaly_score"].(float64)
	totalAmount := result["total_anomalous_amount"].(primitive.Decimal128)

	return &AnomalySummary{
		Count:                count,
		AvgAnomalyScore:      avgScore,
		TotalAnomalousAmount: totalAmount,
	}, nil
}

type SpendingInsight struct {
	Category           string               `json:"category"`
	Merchant           string               `json:"merchant"`
	AvgAmount          primitive.Decimal128 `json:"avg_amount"`
	FrequencyPerMonth  float64              `json:"frequency_per_month"`
	Trend              string               `json:"trend"`
	TrendPercentChange float64              `json:"trend_percent_change"`
	IsRecurring        bool                 `json:"is_recurring"`
	IsEssential        bool                 `json:"is_essential"`
}

func (s *Service) GetTopSpendingInsights(ctx context.Context, userID primitive.ObjectID, limit int) ([]SpendingInsight, error) {
	patterns, err := s.aggregationHelper.GetSpendingTrendsByCategory(ctx, userID, limit)
	if err != nil {
		return nil, err
	}

	insights := make([]SpendingInsight, 0, len(patterns))
	for _, patternMap := range patterns {
		insights = append(insights, SpendingInsight{
			Category:           patternMap["category"].(string),
			Merchant:           patternMap["merchant"].(string),
			AvgAmount:          patternMap["avg_amount"].(primitive.Decimal128),
			FrequencyPerMonth:  patternMap["frequency_per_month"].(float64),
			Trend:              patternMap["trend"].(string),
			TrendPercentChange: patternMap["trend_percent_change"].(float64),
			IsRecurring:        patternMap["is_recurring"].(bool),
			IsEssential:        patternMap["is_essential"].(bool),
		})
	}

	return insights, nil
}

func (s *Service) GetRecurringExpensesSummary(ctx context.Context, userID primitive.ObjectID) (primitive.Decimal128, error) {
	return s.aggregationHelper.GetRecurringExpensesSummary(ctx, userID)
}

type FinancialContext struct {
	Period                string                  `json:"period"`
	Summary               FinancialSummary        `json:"summary"`
	TopSpendingCategories []SpendingCategory      `json:"top_spending_categories"`
	Anomalies             []AnomalyInfo           `json:"anomalies"`
	BudgetStatus          map[string]BudgetStatus `json:"budget_status"`
	Insights              []InsightInfo           `json:"insights"`
	GeneratedAt           time.Time               `json:"generated_at"`
}

func (s *Service) GetAIFinancialContext(ctx context.Context, userID primitive.ObjectID) (*FinancialContext, error) {
	now := time.Now()
	month := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	// Get monthly summary
	summary, err := s.repo.GetMonthlySummary(ctx, userID, month)
	if err != nil {
		return nil, err
	}

	// Get top spending patterns
	patternsData, err := s.aggregationHelper.GetSpendingTrendsByCategory(ctx, userID, 5)
	if err != nil {
		return nil, err
	}

	// Get anomalies
	anomalies, err := s.repo.GetAnomalousTransactions(ctx, userID, 10)
	if err != nil {
		return nil, err
	}

	// Get unread insights
	insights, err := s.repo.GetUnreadInsights(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Build context
	context := &FinancialContext{
		Period: "current_month",
		Summary: FinancialSummary{
			TotalIncome:  summary.Income,
			TotalExpense: summary.Expense,
			NetSavings:   summary.Net,
			SavingsRate:  calculateSavingsRate(summary.Income, summary.Expense),
		},
		TopSpendingCategories: buildSpendingCategoriesFromBsonM(patternsData),
		Anomalies:             buildAnomalies(anomalies),
		BudgetStatus:          make(map[string]BudgetStatus),
		Insights:              buildInsights(insights),
		GeneratedAt:           time.Now(),
	}

	return context, nil
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

func subtractDecimal128(a, b primitive.Decimal128) primitive.Decimal128 {
	// Simplified: In production, use proper Decimal128 arithmetic
	return a
}

func calculateSavingsRate(income, expense primitive.Decimal128) float64 {
	// Simplified: In production, convert Decimal128 to float64 properly
	return 0.0
}

func buildSpendingCategoriesFromBsonM(patterns []bson.M) []SpendingCategory {
	categories := make([]SpendingCategory, 0)
	for _, patternMap := range patterns {
		categories = append(categories, SpendingCategory{
			Name:           patternMap["category"].(string),
			Amount:         patternMap["avg_amount"].(primitive.Decimal128),
			PercentOfTotal: 0.0,
			Trend:          patternMap["trend"].(string),
			VsAverage:      "",
		})
	}
	return categories
}

func buildAnomalies(enrichments []AITransactionEnrichment) []AnomalyInfo {
	anomalies := make([]AnomalyInfo, 0)
	for _, e := range enrichments {
		anomalies = append(anomalies, AnomalyInfo{
			Date:        e.Date,
			Type:        "unusual_spending",
			Description: e.Description,
			Category:    e.MerchantCategory,
			Severity:    "info",
		})
	}
	return anomalies
}

func buildInsights(insights []AIFinancialInsight) []InsightInfo {
	infos := make([]InsightInfo, 0)
	for _, insight := range insights {
		infos = append(infos, InsightInfo{
			Type:             insight.InsightType,
			Message:          insight.Description,
			PotentialSavings: primitive.NewDecimal128(0, 0),
		})
	}
	return infos
}
