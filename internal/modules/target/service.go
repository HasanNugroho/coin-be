package target

import (
	"context"
	"fmt"

	"github.com/HasanNugroho/coin-be/internal/modules/allocation"
	"github.com/HasanNugroho/coin-be/internal/modules/target/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo           *Repository
	allocationRepo *allocation.Repository
}

func NewService(repo *Repository, allocationRepo *allocation.Repository) *Service {
	return &Service{
		repo:           repo,
		allocationRepo: allocationRepo,
	}
}

func (s *Service) CreateTarget(ctx context.Context, userID primitive.ObjectID, req *dto.CreateTargetRequest) (*dto.TargetResponse, error) {
	allocationID, err := primitive.ObjectIDFromHex(req.AllocationID)
	if err != nil {
		return nil, fmt.Errorf("invalid allocation ID: %w", err)
	}

	alloc, err := s.allocationRepo.GetByID(ctx, allocationID)
	if err != nil {
		return nil, fmt.Errorf("allocation not found: %w", err)
	}

	if alloc.UserID != userID {
		return nil, fmt.Errorf("allocation does not belong to user")
	}

	target := &SavingTarget{
		UserID:        userID,
		AllocationID:  allocationID,
		Name:          req.Name,
		TargetAmount:  req.TargetAmount,
		CurrentAmount: alloc.CurrentAmount,
		Deadline:      req.Deadline,
		Status:        TargetStatusActive,
	}

	if err := s.repo.Create(ctx, target); err != nil {
		return nil, fmt.Errorf("failed to create target: %w", err)
	}

	progress := 0.0
	if target.TargetAmount > 0 {
		progress = (target.CurrentAmount / target.TargetAmount) * 100
	}

	return &dto.TargetResponse{
		ID:            target.ID,
		UserID:        target.UserID,
		AllocationID:  target.AllocationID,
		Name:          target.Name,
		TargetAmount:  target.TargetAmount,
		CurrentAmount: target.CurrentAmount,
		Progress:      progress,
		Deadline:      target.Deadline,
		Status:        string(target.Status),
		CreatedAt:     target.CreatedAt,
	}, nil
}

func (s *Service) GetTargets(ctx context.Context, userID primitive.ObjectID) ([]*dto.TargetResponse, error) {
	targets, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get targets: %w", err)
	}

	responses := make([]*dto.TargetResponse, len(targets))
	for i, target := range targets {
		progress := 0.0
		if target.TargetAmount > 0 {
			progress = (target.CurrentAmount / target.TargetAmount) * 100
		}

		responses[i] = &dto.TargetResponse{
			ID:            target.ID,
			UserID:        target.UserID,
			AllocationID:  target.AllocationID,
			Name:          target.Name,
			TargetAmount:  target.TargetAmount,
			CurrentAmount: target.CurrentAmount,
			Progress:      progress,
			Deadline:      target.Deadline,
			Status:        string(target.Status),
			CreatedAt:     target.CreatedAt,
		}
	}

	return responses, nil
}

func (s *Service) GetTargetByID(ctx context.Context, id primitive.ObjectID) (*dto.TargetResponse, error) {
	target, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get target: %w", err)
	}

	progress := 0.0
	if target.TargetAmount > 0 {
		progress = (target.CurrentAmount / target.TargetAmount) * 100
	}

	return &dto.TargetResponse{
		ID:            target.ID,
		UserID:        target.UserID,
		AllocationID:  target.AllocationID,
		Name:          target.Name,
		TargetAmount:  target.TargetAmount,
		CurrentAmount: target.CurrentAmount,
		Progress:      progress,
		Deadline:      target.Deadline,
		Status:        string(target.Status),
		CreatedAt:     target.CreatedAt,
	}, nil
}

func (s *Service) UpdateTarget(ctx context.Context, id primitive.ObjectID, req *dto.UpdateTargetRequest) (*dto.TargetResponse, error) {
	target, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get target: %w", err)
	}

	if req.Name != "" {
		target.Name = req.Name
	}
	if req.TargetAmount > 0 {
		target.TargetAmount = req.TargetAmount
	}
	if !req.Deadline.IsZero() {
		target.Deadline = req.Deadline
	}

	if err := s.repo.Update(ctx, id, target); err != nil {
		return nil, fmt.Errorf("failed to update target: %w", err)
	}

	progress := 0.0
	if target.TargetAmount > 0 {
		progress = (target.CurrentAmount / target.TargetAmount) * 100
	}

	return &dto.TargetResponse{
		ID:            target.ID,
		UserID:        target.UserID,
		AllocationID:  target.AllocationID,
		Name:          target.Name,
		TargetAmount:  target.TargetAmount,
		CurrentAmount: target.CurrentAmount,
		Progress:      progress,
		Deadline:      target.Deadline,
		Status:        string(target.Status),
		CreatedAt:     target.CreatedAt,
	}, nil
}

func (s *Service) DeleteTarget(ctx context.Context, id primitive.ObjectID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete target: %w", err)
	}
	return nil
}

func (s *Service) SyncTargetWithAllocation(ctx context.Context, allocationID primitive.ObjectID) error {
	alloc, err := s.allocationRepo.GetByID(ctx, allocationID)
	if err != nil {
		return fmt.Errorf("allocation not found: %w", err)
	}

	targets, err := s.repo.GetByAllocationID(ctx, allocationID)
	if err != nil {
		return fmt.Errorf("failed to get targets: %w", err)
	}

	for _, target := range targets {
		if target.Status == TargetStatusActive {
			if err := s.repo.UpdateCurrentAmount(ctx, target.ID, alloc.CurrentAmount); err != nil {
				return fmt.Errorf("failed to update target amount: %w", err)
			}

			if alloc.CurrentAmount >= target.TargetAmount {
				if err := s.repo.UpdateStatus(ctx, target.ID, TargetStatusCompleted); err != nil {
					return fmt.Errorf("failed to update target status: %w", err)
				}
			}
		}
	}

	return nil
}
