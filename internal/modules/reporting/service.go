package reporting

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

<<<<<<< Updated upstream
type DashboardService struct {
	repo              *Repository
	aggregationHelper *AggregationHelper
}

func NewService(repo *Repository, aggregationHelper *AggregationHelper) *DashboardService {
	return &DashboardService{
		repo:              repo,
		aggregationHelper: aggregationHelper,
=======
type Service struct {
	repo *Repository
	db   *mongo.Database
}

func NewService(db *mongo.Database) *Service {
	return &Service{
		repo: NewRepository(db),
		db:   db,
>>>>>>> Stashed changes
	}
}

func (s *Service) GetDailyReport(ctx context.Context, userID primitive.ObjectID, date time.Time) (*DailyFinancialReport, error) {
	return s.repo.GetDailyReport(ctx, userID, date)
}

<<<<<<< Updated upstream
func (s *DashboardService) GetDashboardKPIs(ctx context.Context, userID primitive.ObjectID) (*DashboardKPIs, error) {
	now := time.Now()
	month := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
=======
func (s *Service) GenerateDailyReport(ctx context.Context, userID primitive.ObjectID, reportDate time.Time) (*DailyFinancialReport, error) {
	startOfDay := time.Date(reportDate.Year(), reportDate.Month(), reportDate.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.AddDate(0, 0, 1)
>>>>>>> Stashed changes

	txnColl := s.db.Collection("transactions")
	pocketColl := s.db.Collection("pockets")

<<<<<<< Updated upstream
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

func (s *DashboardService) GetDashboardCharts(ctx context.Context, userID primitive.ObjectID) (*DashboardCharts, error) {
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
// REPORTING DashboardService
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

func (s *DashboardService) GetDailyReportSummary(ctx context.Context, userID primitive.ObjectID, reportDate time.Time) (*ReportSummary, error) {
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

func (s *DashboardService) GetMonthlyReportSummary(ctx context.Context, userID primitive.ObjectID, month time.Time) (*ReportSummary, error) {
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

func (s *DashboardService) GetAnomalySummary(ctx context.Context, userID primitive.ObjectID, days int) (*AnomalySummary, error) {
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

func (s *DashboardService) GetTopSpendingInsights(ctx context.Context, userID primitive.ObjectID, limit int) ([]SpendingInsight, error) {
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

func (s *DashboardService) GetRecurringExpensesSummary(ctx context.Context, userID primitive.ObjectID) (primitive.Decimal128, error) {
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

func (s *DashboardService) GetAIFinancialContext(ctx context.Context, userID primitive.ObjectID) (*FinancialContext, error) {
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
=======
	// Fetch all transactions for the day
	txnFilter := bson.M{
		"user_id": userID,
		"date": bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
>>>>>>> Stashed changes
		},
		"deleted_at": bson.M{"$eq": nil},
	}

	cursor, err := txnColl.Find(ctx, txnFilter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []bson.M
	if err = cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}

	// Fetch all active pockets for balance calculation
	pocketFilter := bson.M{
		"user_id":    userID,
		"is_active":  true,
		"deleted_at": bson.M{"$eq": nil},
	}

	pocketCursor, err := pocketColl.Find(ctx, pocketFilter)
	if err != nil {
		return nil, err
	}
	defer pocketCursor.Close(ctx)

	var pockets []bson.M
	if err = pocketCursor.All(ctx, &pockets); err != nil {
		return nil, err
	}

	// Calculate aggregates
	report := &DailyFinancialReport{
		UserID:                      userID,
		ReportDate:                  startOfDay,
		GeneratedAt:                 time.Now(),
		IsFinal:                     false,
		ExpenseByCategory:           []CategoryBreakdown{},
		TransactionsGroupedByPocket: []PocketBreakdown{},
	}

	// Calculate closing balance from current pocket balances
	for _, pocket := range pockets {
		balance, ok := pocket["balance"].(float64)
		if !ok {
			continue
		}
		report.ClosingBalance += balance
	}

	// Aggregate transactions
	categoryMap := make(map[string]*CategoryBreakdown)
	pocketMap := make(map[string]*PocketBreakdown)

	for _, txn := range transactions {
		txnType, _ := txn["type"].(string)
		amount, _ := txn["amount"].(float64)

		switch txnType {
		case "income":
			report.TotalIncome += amount
		case "expense":
			report.TotalExpense += amount
			// Track by category
			if catID, ok := txn["category_id"].(primitive.ObjectID); ok {
				key := catID.Hex()
				if _, exists := categoryMap[key]; !exists {
					categoryMap[key] = &CategoryBreakdown{
						CategoryID: catID,
					}
				}
				categoryMap[key].Amount += amount
				categoryMap[key].TransactionCount++
			}
		case "transfer":
			if _, ok := txn["pocket_to"].(primitive.ObjectID); ok {
				report.TotalTransferIn += amount
			}
			if _, ok := txn["pocket_from"].(primitive.ObjectID); ok {
				report.TotalTransferOut += amount
			}
		}
	}

	// Convert category map to slice
	for _, cat := range categoryMap {
		report.ExpenseByCategory = append(report.ExpenseByCategory, *cat)
	}

	// Convert pocket map to slice
	for _, pocket := range pocketMap {
		report.TransactionsGroupedByPocket = append(report.TransactionsGroupedByPocket, *pocket)
	}

	// Upsert report
	if err := s.repo.UpsertDailyReport(ctx, report); err != nil {
		return nil, err
	}

	return report, nil
}

func (s *Service) GenerateDailySnapshot(ctx context.Context, userID primitive.ObjectID, snapshotDate time.Time) (*DailyFinancialSnapshot, error) {
	pocketColl := s.db.Collection("pockets")

	filter := bson.M{
		"user_id":    userID,
		"is_active":  true,
		"deleted_at": nil,
	}

	cursor, err := pocketColl.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	snapshot := &DailyFinancialSnapshot{
		UserID:         userID,
		SnapshotDate:   snapshotDate,
		PocketBalances: []PocketSnapshot{},
	}

	var pockets []bson.M
	if err = cursor.All(ctx, &pockets); err != nil {
		return nil, err
	}

	for _, pocket := range pockets {
		pocketID, _ := pocket["_id"].(primitive.ObjectID)
		pocketName, _ := pocket["name"].(string)
		pocketType, _ := pocket["type"].(string)
		balance, _ := pocket["balance"].(float64)

		ps := PocketSnapshot{
			PocketID:   pocketID,
			PocketName: pocketName,
			PocketType: pocketType,
			Balance:    balance,
			Currency:   "IDR",
		}

		snapshot.PocketBalances = append(snapshot.PocketBalances, ps)
		snapshot.TotalBalance += balance

		// Free money = main + allocation pockets
		if pocketType == "main" || pocketType == "allocation" {
			snapshot.FreeMoneyTotal += balance
		}
	}

	if err := s.repo.UpsertDailySnapshot(ctx, snapshot); err != nil {
		return nil, err
	}

	return snapshot, nil
}

func (s *Service) GenerateMonthlySummary(ctx context.Context, userID primitive.ObjectID, yearMonth string) (*MonthlyFinancialSummary, error) {
	// Parse year_month to get start and end dates
	parts := len(yearMonth)
	if parts != 7 {
		return nil, fmt.Errorf("invalid year_month format, expected YYYY-MM")
	}

	year := yearMonth[:4]
	month := yearMonth[5:7]

	startDate, err := time.Parse("2006-01-02", fmt.Sprintf("%s-%s-01", year, month))
	if err != nil {
		return nil, err
	}

	endDate := startDate.AddDate(0, 1, 0)

	txnColl := s.db.Collection("transactions")

	filter := bson.M{
		"user_id": userID,
		"date": bson.M{
			"$gte": startDate,
			"$lt":  endDate,
		},
		"deleted_at": bson.M{"$eq": nil},
	}

	cursor, err := txnColl.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []bson.M
	if err = cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}

	summary := &MonthlyFinancialSummary{
		UserID:            userID,
		YearMonth:         yearMonth,
		GeneratedAt:       time.Now(),
		IsFinal:           false,
		ExpenseByCategory: []CategoryBreakdown{},
		IncomeByCategory:  []CategoryBreakdown{},
		PocketSummary:     []PocketSummary{},
	}

	categoryExpenseMap := make(map[string]*CategoryBreakdown)
	categoryIncomeMap := make(map[string]*CategoryBreakdown)
	pocketMap := make(map[string]*PocketSummary)

	for _, txn := range transactions {
		txnType, _ := txn["type"].(string)
		amount, _ := txn["amount"].(float64)

		switch txnType {
		case "income":
			summary.TotalIncome += amount
			if catID, ok := txn["category_id"].(primitive.ObjectID); ok {
				key := catID.Hex()
				if _, exists := categoryIncomeMap[key]; !exists {
					categoryIncomeMap[key] = &CategoryBreakdown{
						CategoryID: catID,
					}
				}
				categoryIncomeMap[key].Amount += amount
				categoryIncomeMap[key].TransactionCount++
			}
		case "expense":
			summary.TotalExpense += amount
			if catID, ok := txn["category_id"].(primitive.ObjectID); ok {
				key := catID.Hex()
				if _, exists := categoryExpenseMap[key]; !exists {
					categoryExpenseMap[key] = &CategoryBreakdown{
						CategoryID: catID,
					}
				}
				categoryExpenseMap[key].Amount += amount
				categoryExpenseMap[key].TransactionCount++
			}
		case "transfer":
			if _, ok := txn["pocket_to"].(primitive.ObjectID); ok {
				summary.TotalTransferIn += amount
			}
			if _, ok := txn["pocket_from"].(primitive.ObjectID); ok {
				summary.TotalTransferOut += amount
			}
		}

		summary.TransactionCount++
	}

	// Convert maps to slices
	for _, cat := range categoryExpenseMap {
		summary.ExpenseByCategory = append(summary.ExpenseByCategory, *cat)
	}
	for _, cat := range categoryIncomeMap {
		summary.IncomeByCategory = append(summary.IncomeByCategory, *cat)
	}
	for _, pocket := range pocketMap {
		summary.PocketSummary = append(summary.PocketSummary, *pocket)
	}

	if err := s.repo.UpsertMonthlySummary(ctx, summary); err != nil {
		return nil, err
	}

	return summary, nil
}

func (s *Service) GetMonthlySummary(ctx context.Context, userID primitive.ObjectID, yearMonth string) (*MonthlyFinancialSummary, error) {
	return s.repo.GetMonthlySummary(ctx, userID, yearMonth)
}

func (s *Service) GetMonthlyTrend(ctx context.Context, userID primitive.ObjectID, months int64) ([]MonthlyFinancialSummary, error) {
	if months > 36 {
		months = 36
	}

	now := time.Now()
	endMonth := fmt.Sprintf("%04d-%02d", now.Year(), now.Month())
	startDate := now.AddDate(0, -int(months), 0)
	startMonth := fmt.Sprintf("%04d-%02d", startDate.Year(), startDate.Month())

	return s.repo.GetMonthlySummariesRange(ctx, userID, startMonth, endMonth, months)
}

func (s *Service) GenerateAIFinancialContext(ctx context.Context, userID primitive.ObjectID) (*AIFinancialContext, error) {
	now := time.Now()
	currentMonth := fmt.Sprintf("%04d-%02d", now.Year(), now.Month())
	thirtyDaysAgo := now.AddDate(0, 0, -30)

	// Get current snapshot
	snapshot, err := s.repo.GetDailySnapshot(ctx, userID, now)
	if err != nil {
		return nil, err
	}

	// Get current month summary
	monthlySummary, err := s.repo.GetMonthlySummary(ctx, userID, currentMonth)
	if err != nil {
		return nil, err
	}

	// Get last 30 days transactions for patterns
	txnColl := s.db.Collection("transactions")
	filter := bson.M{
		"user_id": userID,
		"date": bson.M{
			"$gte": thirtyDaysAgo,
		},
		"deleted_at": bson.M{"$eq": nil},
	}

	cursor, err := txnColl.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []bson.M
	if err = cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}

	aiContext := &AIFinancialContext{
		UserID:      userID,
		ContextDate: now,
		Alerts:      []Alert{},
		Pockets:     []AIPocketData{},
	}

	if snapshot != nil {
		aiContext.CurrentBalance = snapshot.TotalBalance
		aiContext.FreeMoney = snapshot.FreeMoneyTotal
	}

	if monthlySummary != nil {
		aiContext.Last30Days = Last30DaysMetrics{
			TotalIncome:          monthlySummary.TotalIncome,
			TotalExpense:         monthlySummary.TotalExpense,
			NetChange:            monthlySummary.TotalIncome - monthlySummary.TotalExpense,
			TransactionCount:     monthlySummary.TransactionCount,
			TopExpenseCategories: []TopCategory{},
		}

		if len(transactions) > 0 {
			aiContext.Last30Days.AverageDailyExpense = monthlySummary.TotalExpense / 30
		}

		for _, cat := range monthlySummary.ExpenseByCategory {
			if len(aiContext.Last30Days.TopExpenseCategories) < 5 {
				percentage := (cat.Amount / monthlySummary.TotalExpense) * 100
				aiContext.Last30Days.TopExpenseCategories = append(aiContext.Last30Days.TopExpenseCategories, TopCategory{
					CategoryName: cat.CategoryName,
					Amount:       cat.Amount,
					Percentage:   percentage,
				})
			}
		}
	}

	// Calculate year-to-date metrics
	yearStart := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	ytdFilter := bson.M{
		"user_id": userID,
		"date": bson.M{
			"$gte": yearStart,
			"$lt":  now,
		},
		"deleted_at": bson.M{"$eq": nil},
	}

	ytdCursor, err := txnColl.Find(ctx, ytdFilter)
	if err != nil {
		return nil, err
	}
	defer ytdCursor.Close(ctx)

	var ytdTransactions []bson.M
	if err = ytdCursor.All(ctx, &ytdTransactions); err != nil {
		return nil, err
	}

	ytdMetrics := YearToDateMetrics{}
	for _, txn := range ytdTransactions {
		txnType, _ := txn["type"].(string)
		amount, _ := txn["amount"].(float64)

		if txnType == "income" {
			ytdMetrics.TotalIncome += amount
		} else if txnType == "expense" {
			ytdMetrics.TotalExpense += amount
		}
	}

	ytdMetrics.NetChange = ytdMetrics.TotalIncome - ytdMetrics.TotalExpense
	if now.Month() > 1 {
		ytdMetrics.AverageMonthlyExpense = ytdMetrics.TotalExpense / float64(now.Month())
	}

	aiContext.YearToDate = ytdMetrics

	if err := s.repo.UpsertAIFinancialContext(ctx, aiContext); err != nil {
		return nil, err
	}

	return aiContext, nil
}

// GetRealtimeDashboardKPIs retrieves real-time KPIs by querying live collections
func (s *Service) GetRealtimeDashboardKPIs(ctx context.Context, userID primitive.ObjectID, month string) (map[string]interface{}, error) {
	// Parse month
	var startDate, endDate time.Time
	if month == "" {
		now := time.Now()
		month = now.Format("2006-01")
	}

	parsedMonth, err := time.Parse("2006-01", month)
	if err != nil {
		return nil, fmt.Errorf("invalid month format: %w", err)
	}

	startDate = time.Date(parsedMonth.Year(), parsedMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
	endDate = startDate.AddDate(0, 1, 0)

	// Get real-time total balance from pockets
	pocketFilter := bson.M{
		"user_id":    userID,
		"is_active":  true,
		"deleted_at": bson.M{"$eq": nil},
	}

	cursor, err := s.db.Collection("pockets").Find(ctx, pocketFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pockets: %w", err)
	}
	defer cursor.Close(ctx)

	totalBalance := 0.0
	freeMoneyTotal := 0.0

	for cursor.Next(ctx) {
		var pocket struct {
			Balance primitive.Decimal128 `bson:"balance"`
			Type    string               `bson:"type"`
		}
		if err := cursor.Decode(&pocket); err != nil {
			continue
		}

		balanceBigInt, _, err := pocket.Balance.BigInt()
		if err != nil {
			continue
		}
		balanceValue, _ := balanceBigInt.Float64()
		totalBalance += balanceValue

		if pocket.Type == "main" {
			freeMoneyTotal += balanceValue
		}
	}

	// Get real-time income/expense for current month from transactions
	txnFilter := bson.M{
		"user_id": userID,
		"date": bson.M{
			"$gte": startDate,
			"$lt":  endDate,
		},
		"deleted_at": bson.M{"$eq": nil},
	}

	txnCursor, err := s.db.Collection("transactions").Find(ctx, txnFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}
	defer txnCursor.Close(ctx)

	totalIncome := 0.0
	totalExpense := 0.0

	for txnCursor.Next(ctx) {
		var txn struct {
			Type   string               `bson:"type"`
			Amount primitive.Decimal128 `bson:"amount"`
		}
		if err := txnCursor.Decode(&txn); err != nil {
			continue
		}

		amountBigInt, _, err := txn.Amount.BigInt()
		if err != nil {
			continue
		}
		amount, _ := amountBigInt.Float64()

		switch txn.Type {
		case "income":
			totalIncome += amount
		case "expense":
			totalExpense += amount
		}
	}

	return map[string]interface{}{
		"total_balance":               totalBalance,
		"total_income_current_month":  totalIncome,
		"total_expense_current_month": totalExpense,
		"free_money_total":            freeMoneyTotal,
		"net_change_current_month":    totalIncome - totalExpense,
	}, nil
}

// GetRealtimePocketDistribution retrieves real-time pocket distribution
func (s *Service) GetRealtimePocketDistribution(ctx context.Context, userID primitive.ObjectID) ([]map[string]interface{}, error) {
	pocketFilter := bson.M{
		"user_id":    userID,
		"is_active":  true,
		"deleted_at": bson.M{"$eq": nil},
	}

	cursor, err := s.db.Collection("pockets").Find(ctx, pocketFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pockets: %w", err)
	}
	defer cursor.Close(ctx)

	var pockets []struct {
		ID      primitive.ObjectID   `bson:"_id"`
		Name    string               `bson:"name"`
		Type    string               `bson:"type"`
		Balance primitive.Decimal128 `bson:"balance"`
	}

	if err := cursor.All(ctx, &pockets); err != nil {
		return nil, fmt.Errorf("failed to decode pockets: %w", err)
	}

	// Calculate total balance
	totalBalance := 0.0
	for _, p := range pockets {
		balanceBigInt, _, err := p.Balance.BigInt()
		if err != nil {
			continue
		}
		balance, _ := balanceBigInt.Float64()
		totalBalance += balance
	}

	// Build response
	result := make([]map[string]interface{}, 0, len(pockets))
	for _, p := range pockets {
		balanceBigInt, _, err := p.Balance.BigInt()
		if err != nil {
			continue
		}
		balance, _ := balanceBigInt.Float64()
		percentage := 0.0
		if totalBalance > 0 {
			percentage = (balance / totalBalance) * 100
		}

		result = append(result, map[string]interface{}{
			"pocket_id":   p.ID.Hex(),
			"pocket_name": p.Name,
			"pocket_type": p.Type,
			"balance":     balance,
			"percentage":  percentage,
		})
	}

	return result, nil
}

// GetRealtimeExpenseByCategory retrieves real-time expense distribution by category
func (s *Service) GetRealtimeExpenseByCategory(ctx context.Context, userID primitive.ObjectID, month string) ([]map[string]interface{}, error) {
	// Parse month
	var startDate, endDate time.Time
	if month == "" {
		now := time.Now()
		month = now.Format("2006-01")
	}

	parsedMonth, err := time.Parse("2006-01", month)
	if err != nil {
		return nil, fmt.Errorf("invalid month format: %w", err)
	}

	startDate = time.Date(parsedMonth.Year(), parsedMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
	endDate = startDate.AddDate(0, 1, 0)

	// Aggregate expenses by category
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"user_id": userID,
				"type":    "expense",
				"date": bson.M{
					"$gte": startDate,
					"$lt":  endDate,
				},
				"deleted_at": bson.M{"$eq": nil},
			},
		},
		{
			"$group": bson.M{
				"_id": "$category_id",
				"total_amount": bson.M{
					"$sum": bson.M{"$toDouble": "$amount"},
				},
				"transaction_count": bson.M{"$sum": 1},
			},
		},
		{
			"$sort": bson.M{"total_amount": -1},
		},
	}

	cursor, err := s.db.Collection("transactions").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate expenses: %w", err)
	}
	defer cursor.Close(ctx)

	type CategoryAggregate struct {
		CategoryID       primitive.ObjectID `bson:"_id"`
		TotalAmount      float64            `bson:"total_amount"`
		TransactionCount int32              `bson:"transaction_count"`
	}

	var aggregates []CategoryAggregate
	if err := cursor.All(ctx, &aggregates); err != nil {
		return nil, fmt.Errorf("failed to decode aggregates: %w", err)
	}

	// Calculate total expense
	totalExpense := 0.0
	for _, agg := range aggregates {
		totalExpense += agg.TotalAmount
	}

	// Fetch category names
	categoryIDs := make([]primitive.ObjectID, 0, len(aggregates))
	for _, agg := range aggregates {
		categoryIDs = append(categoryIDs, agg.CategoryID)
	}

	categoryMap := make(map[primitive.ObjectID]string)
	if len(categoryIDs) > 0 {
		categoryCursor, err := s.db.Collection("user_categories").Find(ctx, bson.M{
			"_id": bson.M{"$in": categoryIDs},
		})
		if err == nil {
			defer categoryCursor.Close(ctx)

			for categoryCursor.Next(ctx) {
				var cat struct {
					ID   primitive.ObjectID `bson:"_id"`
					Name string             `bson:"name"`
				}
				if err := categoryCursor.Decode(&cat); err == nil {
					categoryMap[cat.ID] = cat.Name
				}
			}
		}
	}

	// Build response
	result := make([]map[string]interface{}, 0, len(aggregates))
	for _, agg := range aggregates {
		percentage := 0.0
		if totalExpense > 0 {
			percentage = (agg.TotalAmount / totalExpense) * 100
		}

		categoryName := categoryMap[agg.CategoryID]
		if categoryName == "" {
			categoryName = "Uncategorized"
		}

		result = append(result, map[string]interface{}{
			"category_id":       agg.CategoryID.Hex(),
			"category_name":     categoryName,
			"amount":            agg.TotalAmount,
			"percentage":        percentage,
			"transaction_count": agg.TransactionCount,
		})
	}

	return result, nil
}

func (s *Service) GetAIFinancialContext(ctx context.Context, userID primitive.ObjectID) (*AIFinancialContext, error) {
	return s.repo.GetAIFinancialContext(ctx, userID)
}
