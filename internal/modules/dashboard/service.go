package dashboard

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{
		repo: r,
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

	startOfPeriod, startOfToday := timeRange.ToDuration()

	historicalIncome, historicalExpense, _, err := s.repo.GetHistoricalSummary(ctx, userObjID, startOfPeriod, startOfToday)
	if err != nil {
		return nil, err
	}

	liveIncome, liveExpense, _, err := s.repo.GetLiveDeltaSummary(ctx, userObjID, startOfToday)
	if err != nil {
		return nil, err
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

	startOfPeriod, startOfToday := timeRange.ToDuration()

	cashFlowTrend, err := s.repo.GetDailyCashFlowTrend(ctx, userObjID, startOfPeriod, startOfToday)
	if err != nil {
		return nil, err
	}

	historicalIncome, historicalExpense, historicalCategories, err := s.repo.GetHistoricalSummary(ctx, userObjID, startOfPeriod, startOfToday)
	if err != nil {
		return nil, err
	}

	liveIncome, liveExpense, liveCategories, err := s.repo.GetLiveDeltaSummary(ctx, userObjID, startOfToday)
	if err != nil {
		return nil, err
	}

	totalIncome := historicalIncome + liveIncome
	totalExpense := historicalExpense + liveExpense

	categoryMap := make(map[string]*CategoryBreakdown)
	for _, cat := range historicalCategories {
		key := cat.Type + "_"
		if cat.CategoryID != nil {
			key += cat.CategoryID.Hex()
		} else {
			key += "uncategorized"
		}
		categoryMap[key] = &CategoryBreakdown{
			CategoryID:   cat.CategoryID,
			CategoryName: cat.CategoryName,
			Type:         cat.Type,
			Amount:       cat.Amount,
		}
	}

	for _, cat := range liveCategories {
		key := cat.Type + "_"
		if cat.CategoryID != nil {
			key += cat.CategoryID.Hex()
		} else {
			key += "uncategorized"
		}

		if existing, ok := categoryMap[key]; ok {
			existing.Amount += cat.Amount
		} else {
			categoryMap[key] = &CategoryBreakdown{
				CategoryID:   cat.CategoryID,
				CategoryName: cat.CategoryName,
				Type:         cat.Type,
				Amount:       cat.Amount,
			}
		}
	}

	incomeBreakdown := []CategoryChartData{}
	expenseBreakdown := []CategoryChartData{}

	for _, cat := range categoryMap {
		categoryID := ""
		if cat.CategoryID != nil {
			categoryID = cat.CategoryID.Hex()
		}

		percentage := 0.0
		if cat.Type == "income" && totalIncome > 0 {
			percentage = (cat.Amount / totalIncome) * 100
			incomeBreakdown = append(incomeBreakdown, CategoryChartData{
				CategoryID:   categoryID,
				CategoryName: cat.CategoryName,
				Amount:       cat.Amount,
				Percentage:   percentage,
			})
		} else if cat.Type == "expense" && totalExpense > 0 {
			percentage = (cat.Amount / totalExpense) * 100
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

func (s *Service) GenerateDailySummary(ctx context.Context, userID primitive.ObjectID, date time.Time) error {
	return s.repo.GenerateDailySummaryForDate(ctx, userID, date)
}

func (s *Service) GenerateDailySummariesForAllUsers(ctx context.Context, date time.Time) error {
	userIDs, err := s.repo.GetAllUsersWithTransactions(ctx, date)
	if err != nil {
		return err
	}

	for _, userID := range userIDs {
		if err := s.repo.GenerateDailySummaryForDate(ctx, userID, date); err != nil {
			continue
		}
	}

	return nil
}
