package dashboard

import (
	"context"
	"errors"
	"time"

	"github.com/HasanNugroho/coin-be/internal/modules/daily_summary"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo             *Repository
	dailySummaryRepo *daily_summary.Repository
}

func NewService(r *Repository, dsr *daily_summary.Repository) *Service {
	return &Service{
		repo:             r,
		dailySummaryRepo: dsr,
	}
}

type TimeRange string

const (
	TimeRange7Days  TimeRange = "7d"
	TimeRange1Month TimeRange = "1m"
	TimeRange3Month TimeRange = "3m"
)

func (t TimeRange) ToDuration() (time.Time, time.Time) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	switch t {
	case TimeRange7Days:
		return today.AddDate(0, 0, -7), today
	case TimeRange1Month:
		// calendar: 1st of current month
		return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()), today
	case TimeRange3Month:
		// calendar: 1st of the month, 3 months ago
		start := today.AddDate(0, -3, 0)
		return time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, now.Location()), today
	default:
		// rolling 30 days
		return today.AddDate(0, 0, -30), today
	}
}

func (s *Service) GetDashboardSummary(ctx context.Context, userID string, timeRange TimeRange) (*DashboardSummary, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	totalNetWorth, err := s.repo.GetTotalNetWorth(ctx, userObjID)
	if err != nil {
		return nil, err
	}

	if timeRange == "" {
		timeRange = TimeRange1Month
	}

	startDate, today := timeRange.ToDuration()

	summaries, err := s.dailySummaryRepo.GetDailySummariesByDateRange(ctx, userObjID, startDate, today)
	if err != nil {
		return nil, err
	}

	liveIncome, liveExpense, _, err := s.repo.GetLiveDeltaSummary(ctx, userObjID, today)
	if err != nil {
		return nil, err
	}

	historicalIncome := 0.0
	historicalExpense := 0.0
	for _, s := range summaries {
		historicalIncome += s.TotalIncome
		historicalExpense += s.TotalExpense
	}

	periodIncome := historicalIncome + liveIncome
	periodExpense := historicalExpense + liveExpense
	periodNet := periodIncome - periodExpense

	return &DashboardSummary{
		TotalNetWorth: totalNetWorth,
		PeriodIncome:  periodIncome,
		PeriodExpense: periodExpense,
		PeriodNet:     periodNet,
		TimeRange:     timeRange,
	}, nil
}

func (s *Service) GetDashboardCharts(ctx context.Context, userID string, timeRange TimeRange) (*DashboardCharts, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	if timeRange == "" {
		timeRange = TimeRange1Month
	}

	startDate, today := timeRange.ToDuration()

	// 1. Fetch historical summaries and live delta
	summaries, err := s.dailySummaryRepo.GetDailySummariesByDateRange(ctx, userObjID, startDate, today)
	if err != nil {
		return nil, err
	}

	liveIncome, liveExpense, liveCategories, err := s.repo.GetLiveDeltaSummary(ctx, userObjID, today)
	if err != nil {
		return nil, err
	}

	// 2. Map summaries and collect IDs for name population
	summaryMap := make(map[string]*daily_summary.DailySummary)
	historicalIncome := 0.0
	historicalExpense := 0.0
	categoryMap := make(map[string]*daily_summary.CategoryBreakdown)

	catIDs := make(map[primitive.ObjectID]bool)

	for _, sm := range summaries {
		dateStr := sm.Date.Format("2006-01-02")
		summaryMap[dateStr] = sm
		historicalIncome += sm.TotalIncome
		historicalExpense += sm.TotalExpense

		for _, cat := range sm.CategoryBreakdown {
			if cat.CategoryID != nil {
				catIDs[*cat.CategoryID] = true
			}
			key := cat.Type + "_"
			if cat.CategoryID != nil {
				key += cat.CategoryID.Hex()
			} else {
				key += "uncategorized"
			}

			if existing, ok := categoryMap[key]; ok {
				existing.Amount += cat.Amount
			} else {
				categoryMap[key] = &daily_summary.CategoryBreakdown{
					CategoryID: cat.CategoryID,
					Type:       cat.Type,
					Amount:     cat.Amount,
				}
			}
		}
	}

	// 3. Populate Category Names in bulk
	if len(catIDs) > 0 {
		ids := make([]primitive.ObjectID, 0, len(catIDs))
		for id := range catIDs {
			ids = append(ids, id)
		}
		names, _ := s.repo.GetCategoryNames(ctx, ids) // We should implement this or use a simple map
		for _, cat := range categoryMap {
			if cat.CategoryID != nil {
				if name, ok := names[*cat.CategoryID]; ok {
					cat.CategoryName = name
				} else {
					cat.CategoryName = "Unknown"
				}
			} else {
				cat.CategoryName = "Uncategorized"
			}
		}
	} else {
		for _, cat := range categoryMap {
			if cat.CategoryID == nil {
				cat.CategoryName = "Uncategorized"
			}
		}
	}

	// 4. Merge live categories into categoryMap
	for _, cat := range liveCategories {
		key := cat.Type + "_"
		if cat.CategoryID != nil {
			key += cat.CategoryID.Hex()
		} else {
			key += "uncategorized"
		}

		if existing, ok := categoryMap[key]; ok {
			existing.Amount += cat.Amount
			// Live category already has name from repo lookup
			if existing.CategoryName == "" || existing.CategoryName == "Unknown" {
				existing.CategoryName = cat.CategoryName
			}
		} else {
			categoryMap[key] = &daily_summary.CategoryBreakdown{
				CategoryID:   cat.CategoryID,
				CategoryName: cat.CategoryName,
				Type:         cat.Type,
				Amount:       cat.Amount,
			}
		}
	}

	totalIncome := historicalIncome + liveIncome
	totalExpense := historicalExpense + liveExpense

	// 5. Build CashFlowTrend
	cashFlowTrend := []ChartDataPoint{}
	for d := startDate; !d.After(today); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		point := ChartDataPoint{Date: dateStr}

		if d.Equal(today) {
			point.Income = liveIncome
			point.Expense = liveExpense
		} else if sm, ok := summaryMap[dateStr]; ok {
			point.Income = sm.TotalIncome
			point.Expense = sm.TotalExpense
		}

		cashFlowTrend = append(cashFlowTrend, point)
	}

	// 6. Build breakdowns
	incomeBreakdown := []CategoryChartData{}
	expenseBreakdown := []CategoryChartData{}

	for _, cat := range categoryMap {
		categoryID := ""
		if cat.CategoryID != nil {
			categoryID = cat.CategoryID.Hex()
		}

		if cat.Type == "income" {
			percentage := 0.0
			if totalIncome > 0 {
				percentage = (cat.Amount / totalIncome) * 100
			}
			incomeBreakdown = append(incomeBreakdown, CategoryChartData{
				CategoryID:   categoryID,
				CategoryName: cat.CategoryName,
				Amount:       cat.Amount,
				Percentage:   percentage,
			})
		} else if cat.Type == "expense" {
			percentage := 0.0
			if totalExpense > 0 {
				percentage = (cat.Amount / totalExpense) * 100
			}
			expenseBreakdown = append(expenseBreakdown, CategoryChartData{
				CategoryID:   categoryID,
				CategoryName: cat.CategoryName,
				Amount:       cat.Amount,
				Percentage:   percentage,
			})
		}
	}

	return &DashboardCharts{
		CashFlowTrend:    cashFlowTrend,
		IncomeBreakdown:  incomeBreakdown,
		ExpenseBreakdown: expenseBreakdown,
	}, nil
}
