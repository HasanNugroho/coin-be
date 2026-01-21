package allocation

import (
	"context"
	"fmt"

	"github.com/HasanNugroho/coin-be/internal/modules/allocation/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateAllocation(ctx context.Context, userID primitive.ObjectID, req *dto.CreateAllocationRequest) (*dto.AllocationResponse, error) {
	allocation := &Allocation{
		UserID:        userID,
		Name:          req.Name,
		Priority:      req.Priority,
		Percentage:    req.Percentage,
		CurrentAmount: 0,
		TargetAmount:  req.TargetAmount,
		IsActive:      true,
	}

	if err := s.repo.Create(ctx, allocation); err != nil {
		return nil, fmt.Errorf("failed to create allocation: %w", err)
	}

	return &dto.AllocationResponse{
		ID:            allocation.ID,
		UserID:        allocation.UserID,
		Name:          allocation.Name,
		Priority:      allocation.Priority,
		Percentage:    allocation.Percentage,
		CurrentAmount: allocation.CurrentAmount,
		TargetAmount:  allocation.TargetAmount,
		IsActive:      allocation.IsActive,
		CreatedAt:     allocation.CreatedAt,
	}, nil
}

func (s *Service) GetAllocations(ctx context.Context, userID primitive.ObjectID) ([]*dto.AllocationResponse, error) {
	allocations, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get allocations: %w", err)
	}

	responses := make([]*dto.AllocationResponse, len(allocations))
	for i, alloc := range allocations {
		responses[i] = &dto.AllocationResponse{
			ID:            alloc.ID,
			UserID:        alloc.UserID,
			Name:          alloc.Name,
			Priority:      alloc.Priority,
			Percentage:    alloc.Percentage,
			CurrentAmount: alloc.CurrentAmount,
			TargetAmount:  alloc.TargetAmount,
			IsActive:      alloc.IsActive,
			CreatedAt:     alloc.CreatedAt,
		}
	}

	return responses, nil
}

func (s *Service) GetAllocationByID(ctx context.Context, id primitive.ObjectID) (*dto.AllocationResponse, error) {
	allocation, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get allocation: %w", err)
	}

	return &dto.AllocationResponse{
		ID:            allocation.ID,
		UserID:        allocation.UserID,
		Name:          allocation.Name,
		Priority:      allocation.Priority,
		Percentage:    allocation.Percentage,
		CurrentAmount: allocation.CurrentAmount,
		TargetAmount:  allocation.TargetAmount,
		IsActive:      allocation.IsActive,
		CreatedAt:     allocation.CreatedAt,
	}, nil
}

func (s *Service) UpdateAllocation(ctx context.Context, id primitive.ObjectID, req *dto.UpdateAllocationRequest) (*dto.AllocationResponse, error) {
	allocation, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get allocation: %w", err)
	}

	if req.Name != "" {
		allocation.Name = req.Name
	}
	if req.Priority > 0 {
		allocation.Priority = req.Priority
	}
	if req.Percentage > 0 {
		allocation.Percentage = req.Percentage
	}
	if req.TargetAmount != nil {
		allocation.TargetAmount = req.TargetAmount
	}
	if req.IsActive != nil {
		allocation.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, id, allocation); err != nil {
		return nil, fmt.Errorf("failed to update allocation: %w", err)
	}

	return &dto.AllocationResponse{
		ID:            allocation.ID,
		UserID:        allocation.UserID,
		Name:          allocation.Name,
		Priority:      allocation.Priority,
		Percentage:    allocation.Percentage,
		CurrentAmount: allocation.CurrentAmount,
		TargetAmount:  allocation.TargetAmount,
		IsActive:      allocation.IsActive,
		CreatedAt:     allocation.CreatedAt,
	}, nil
}

func (s *Service) DeleteAllocation(ctx context.Context, id primitive.ObjectID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete allocation: %w", err)
	}
	return nil
}

func (s *Service) GetAllocationLogs(ctx context.Context, userID primitive.ObjectID, limit, skip int64) ([]*dto.AllocationLogResponse, error) {
	logs, err := s.repo.GetLogsByUserID(ctx, userID, limit, skip)
	if err != nil {
		return nil, fmt.Errorf("failed to get allocation logs: %w", err)
	}

	responses := make([]*dto.AllocationLogResponse, len(logs))
	for i, log := range logs {
		responses[i] = &dto.AllocationLogResponse{
			ID:              log.ID,
			UserID:          log.UserID,
			AllocationID:    log.AllocationID,
			TransactionID:   log.TransactionID,
			IncomeAmount:    log.IncomeAmount,
			AllocatedAmount: log.AllocatedAmount,
			Percentage:      log.Percentage,
			Priority:        log.Priority,
			CreatedAt:       log.CreatedAt,
		}
	}

	return responses, nil
}

func (s *Service) GetAllocationLogsByAllocationID(ctx context.Context, allocationID primitive.ObjectID, limit, skip int64) ([]*dto.AllocationLogResponse, error) {
	logs, err := s.repo.GetLogsByAllocationID(ctx, allocationID, limit, skip)
	if err != nil {
		return nil, fmt.Errorf("failed to get allocation logs: %w", err)
	}

	responses := make([]*dto.AllocationLogResponse, len(logs))
	for i, log := range logs {
		responses[i] = &dto.AllocationLogResponse{
			ID:              log.ID,
			UserID:          log.UserID,
			AllocationID:    log.AllocationID,
			TransactionID:   log.TransactionID,
			IncomeAmount:    log.IncomeAmount,
			AllocatedAmount: log.AllocatedAmount,
			Percentage:      log.Percentage,
			Priority:        log.Priority,
			CreatedAt:       log.CreatedAt,
		}
	}

	return responses, nil
}
