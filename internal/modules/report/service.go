package report

import (
	"context"
	"fmt"
	"time"

	"github.com/HasanNugroho/coin-be/internal/modules/allocation"
	"github.com/HasanNugroho/coin-be/internal/modules/category"
	"github.com/HasanNugroho/coin-be/internal/modules/target"
	"github.com/HasanNugroho/coin-be/internal/modules/transaction"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	transactionRepo *transaction.Repository
	allocationRepo  *allocation.Repository
	categoryRepo    *category.Repository
	targetRepo      *target.Repository
}

func NewService(
	transactionRepo *transaction.Repository,
	allocationRepo *allocation.Repository,
	categoryRepo *category.Repository,
	targetRepo *target.Repository,
) *Service {
	return &Service{
		transactionRepo: transactionRepo,
		allocationRepo:  allocationRepo,
		categoryRepo:    categoryRepo,
		targetRepo:      targetRepo,
	}
}

func (s *Service) GetDashboardSummary(ctx context.Context, userID primitive.ObjectID) (*DashboardSummary, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)

	totalAllocated, err := s.allocationRepo.GetTotalCurrentAmount(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get total allocated: %w", err)
	}

	incomeThisMonth, err := s.transactionRepo.GetTotalByUserAndType(ctx, userID, transaction.TransactionTypeIncome, startOfMonth, endOfMonth)
	if err != nil {
		return nil, fmt.Errorf("failed to get income: %w", err)
	}

	expenseThisMonth, err := s.transactionRepo.GetTotalByUserAndType(ctx, userID, transaction.TransactionTypeExpense, startOfMonth, endOfMonth)
	if err != nil {
		return nil, fmt.Errorf("failed to get expense: %w", err)
	}

	return &DashboardSummary{
		TotalBalance:       totalAllocated,
		FreeCash:           0,
		IncomeThisMonth:    incomeThisMonth,
		ExpenseThisMonth:   expenseThisMonth,
		RemainingThisMonth: incomeThisMonth - expenseThisMonth,
	}, nil
}

func (s *Service) GetIncomeReport(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) (*IncomeReport, error) {
	totalIncome, err := s.transactionRepo.GetTotalByUserAndType(ctx, userID, transaction.TransactionTypeIncome, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get total income: %w", err)
	}

	categories, err := s.categoryRepo.GetByUserIDAndType(ctx, userID, category.CategoryTypeIncome)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	var byCategory []CategoryBreakdown
	for _, cat := range categories {
		transactions, err := s.transactionRepo.GetByCategory(ctx, userID, cat.ID, startDate, endDate)
		if err != nil {
			continue
		}

		var total float64
		for _, txn := range transactions {
			total += txn.Amount
		}

		if total > 0 {
			percentage := 0.0
			if totalIncome > 0 {
				percentage = (total / totalIncome) * 100
			}

			byCategory = append(byCategory, CategoryBreakdown{
				CategoryID:   cat.ID,
				CategoryName: cat.Name,
				Amount:       total,
				Percentage:   percentage,
				Count:        len(transactions),
			})
		}
	}

	return &IncomeReport{
		TotalIncome: totalIncome,
		ByCategory:  byCategory,
		ByMonth:     []MonthlyBreakdown{},
	}, nil
}

func (s *Service) GetExpenseReport(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) (*ExpenseReport, error) {
	totalExpense, err := s.transactionRepo.GetTotalByUserAndType(ctx, userID, transaction.TransactionTypeExpense, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get total expense: %w", err)
	}

	categories, err := s.categoryRepo.GetByUserIDAndType(ctx, userID, category.CategoryTypeExpense)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	var byCategory []CategoryBreakdown
	for _, cat := range categories {
		transactions, err := s.transactionRepo.GetByCategory(ctx, userID, cat.ID, startDate, endDate)
		if err != nil {
			continue
		}

		var total float64
		for _, txn := range transactions {
			total += txn.Amount
		}

		if total > 0 {
			percentage := 0.0
			if totalExpense > 0 {
				percentage = (total / totalExpense) * 100
			}

			byCategory = append(byCategory, CategoryBreakdown{
				CategoryID:   cat.ID,
				CategoryName: cat.Name,
				Amount:       total,
				Percentage:   percentage,
				Count:        len(transactions),
			})
		}
	}

	allocations, err := s.allocationRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get allocations: %w", err)
	}

	var byAllocation []AllocationBreakdown
	for _, alloc := range allocations {
		transactions, err := s.transactionRepo.GetByAllocation(ctx, userID, alloc.ID, startDate, endDate)
		if err != nil {
			continue
		}

		var total float64
		for _, txn := range transactions {
			total += txn.Amount
		}

		if total > 0 {
			percentage := 0.0
			if totalExpense > 0 {
				percentage = (total / totalExpense) * 100
			}

			byAllocation = append(byAllocation, AllocationBreakdown{
				AllocationID:   alloc.ID,
				AllocationName: alloc.Name,
				Amount:         total,
				Percentage:     percentage,
				Count:          len(transactions),
			})
		}
	}

	return &ExpenseReport{
		TotalExpense: totalExpense,
		ByCategory:   byCategory,
		ByAllocation: byAllocation,
		ByMonth:      []MonthlyBreakdown{},
	}, nil
}

func (s *Service) GetAllocationReport(ctx context.Context, userID primitive.ObjectID) ([]*AllocationReport, error) {
	allocations, err := s.allocationRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get allocations: %w", err)
	}

	var reports []*AllocationReport
	for _, alloc := range allocations {
		logs, err := s.allocationRepo.GetLogsByAllocationID(ctx, alloc.ID, 1000, 0)
		if err != nil {
			continue
		}

		var totalAllocated float64
		for _, log := range logs {
			totalAllocated += log.AllocatedAmount
		}

		progress := 0.0
		if alloc.TargetAmount != nil && *alloc.TargetAmount > 0 {
			progress = (alloc.CurrentAmount / *alloc.TargetAmount) * 100
		}

		reports = append(reports, &AllocationReport{
			AllocationID:      alloc.ID,
			AllocationName:    alloc.Name,
			CurrentBalance:    alloc.CurrentAmount,
			TargetAmount:      alloc.TargetAmount,
			Progress:          progress,
			TotalAllocated:    totalAllocated,
			TotalSpent:        totalAllocated - alloc.CurrentAmount,
			DistributionCount: len(logs),
		})
	}

	return reports, nil
}

func (s *Service) GetTargetProgress(ctx context.Context, userID primitive.ObjectID) ([]*TargetProgress, error) {
	targets, err := s.targetRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get targets: %w", err)
	}

	var progress []*TargetProgress
	now := time.Now()

	for _, tgt := range targets {
		progressPct := 0.0
		if tgt.TargetAmount > 0 {
			progressPct = (tgt.CurrentAmount / tgt.TargetAmount) * 100
		}

		daysRemaining := int(tgt.Deadline.Sub(now).Hours() / 24)
		if daysRemaining < 0 {
			daysRemaining = 0
		}

		progress = append(progress, &TargetProgress{
			TargetID:      tgt.ID,
			TargetName:    tgt.Name,
			TargetAmount:  tgt.TargetAmount,
			CurrentAmount: tgt.CurrentAmount,
			Progress:      progressPct,
			Status:        string(tgt.Status),
			DaysRemaining: daysRemaining,
		})
	}

	return progress, nil
}
