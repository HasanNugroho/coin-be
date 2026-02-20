package admin_dashboard

import (
	"context"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) GetAdminDashboardSummary(ctx context.Context, startDate, endDate time.Time) (*AdminDashboardSummary, error) {
	totalUsers, err := s.repo.GetTotalUsersCount(ctx)
	if err != nil {
		return nil, err
	}

	activeUsers, err := s.repo.GetActiveUsersCount(ctx)
	if err != nil {
		return nil, err
	}

	totalTransactions, err := s.repo.GetTotalTransactionsCount(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	totalVolume, err := s.repo.GetTotalTransactionVolume(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	userGrowth, err := s.repo.GetUserGrowth(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return &AdminDashboardSummary{
		TotalUsers:        totalUsers,
		ActiveUsers:       activeUsers,
		TotalTransactions: totalTransactions,
		TotalVolume:       totalVolume,
		UserGrowth:        userGrowth,
	}, nil
}
