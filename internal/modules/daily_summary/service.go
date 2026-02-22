package daily_summary

import (
	"context"
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

func (s *Service) GenerateDailySummary(ctx context.Context, userID primitive.ObjectID, date time.Time) error {
	return s.repo.GenerateDailySummaryForDate(ctx, userID, date)
}

func (s *Service) GenerateDailySummariesForAllUsers(ctx context.Context, date time.Time) error {
	return s.repo.GenerateDailySummariesFromTo(ctx, date)
}

func (s *Service) SyncDailySummaries(ctx context.Context, startDate time.Time) error {
	// 1. Delete first (Optional, GenerateDailySummariesFromTo also deletes, but let's keep it for safety if we want to be sure everything is wiped)
	if err := s.repo.DeleteDailySummariesByDateRange(ctx, startDate); err != nil {
		return err
	}

	// 2. Simply call the optimized batch generator
	return s.repo.GenerateDailySummariesFromTo(ctx, startDate)
}
