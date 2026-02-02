package allocation

import (
	"context"
	"errors"

	"github.com/HasanNugroho/coin-be/internal/modules/allocation/dto"
	"github.com/HasanNugroho/coin-be/internal/modules/pocket"
	"github.com/HasanNugroho/coin-be/internal/modules/user_platform"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo             *Repository
	pocketRepo       *pocket.Repository
	userPlatformRepo *user_platform.UserPlatformRepository
}

func NewService(r *Repository, pr *pocket.Repository, upr *user_platform.UserPlatformRepository) *Service {
	return &Service{
		repo:             r,
		pocketRepo:       pr,
		userPlatformRepo: upr,
	}
}

func (s *Service) CreateAllocation(ctx context.Context, userID string, req *dto.CreateAllocationRequest) (*Allocation, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	if !IsValidAllocationType(req.AllocationType) {
		return nil, errors.New("invalid allocation type")
	}

	if !IsValidPriority(req.Priority) {
		return nil, errors.New("invalid priority")
	}

	// Validate that at least one target is provided
	if req.PocketID == "" && req.UserPlatformID == "" {
		return nil, errors.New("either pocket_id or user_platform_id must be provided")
	}

	var pocketID *primitive.ObjectID
	var userPlatformID *primitive.ObjectID

	// Validate pocket if provided
	if req.PocketID != "" {
		pocketObjID, err := primitive.ObjectIDFromHex(req.PocketID)
		if err != nil {
			return nil, errors.New("invalid pocket id")
		}

		pocket, err := s.pocketRepo.GetPocketByID(ctx, pocketObjID)
		if err != nil {
			return nil, errors.New("pocket not found")
		}

		if pocket.UserID != userObjID {
			return nil, errors.New("unauthorized: pocket does not belong to user")
		}

		if !pocket.IsActive {
			return nil, errors.New("pocket is not active")
		}

		pocketID = &pocketObjID
	}

	// Validate user platform if provided
	if req.UserPlatformID != "" {
		userPlatformObjID, err := primitive.ObjectIDFromHex(req.UserPlatformID)
		if err != nil {
			return nil, errors.New("invalid user platform id")
		}

		userPlatform, err := s.userPlatformRepo.GetUserPlatformByID(ctx, userPlatformObjID)
		if err != nil {
			return nil, errors.New("user platform not found")
		}

		if userPlatform.UserID != userObjID {
			return nil, errors.New("unauthorized: user platform does not belong to user")
		}

		if !userPlatform.IsActive {
			return nil, errors.New("user platform is not active")
		}

		userPlatformID = &userPlatformObjID
	}

	// Validate nominal based on allocation type
	if req.AllocationType == string(TypePercentage) && req.Nominal > 100 {
		return nil, errors.New("percentage cannot exceed 100")
	}

	allocation := &Allocation{
		UserID:         userObjID,
		PocketID:       pocketID,
		UserPlatformID: userPlatformID,
		Priority:       req.Priority,
		AllocationType: req.AllocationType,
		Nominal:        req.Nominal,
		IsActive:       true,
	}

	err = s.repo.CreateAllocation(ctx, allocation)
	if err != nil {
		return nil, err
	}

	return allocation, nil
}

func (s *Service) GetAllocationByID(ctx context.Context, userID string, allocationID string) (*Allocation, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	allocationObjID, err := primitive.ObjectIDFromHex(allocationID)
	if err != nil {
		return nil, errors.New("invalid allocation id")
	}

	allocation, err := s.repo.GetAllocationByID(ctx, allocationObjID)
	if err != nil {
		return nil, err
	}

	if allocation.UserID != userObjID {
		return nil, errors.New("unauthorized")
	}

	return allocation, nil
}

func (s *Service) ListAllocations(ctx context.Context, userID string) ([]*Allocation, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	return s.repo.GetAllocationsByUserID(ctx, userObjID)
}

func (s *Service) GetActiveAllocations(ctx context.Context, userID string) ([]*Allocation, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	return s.repo.GetActiveAllocationsByUserID(ctx, userObjID)
}

func (s *Service) UpdateAllocation(ctx context.Context, userID string, allocationID string, req *dto.UpdateAllocationRequest) (*Allocation, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	allocationObjID, err := primitive.ObjectIDFromHex(allocationID)
	if err != nil {
		return nil, errors.New("invalid allocation id")
	}

	allocation, err := s.repo.GetAllocationByID(ctx, allocationObjID)
	if err != nil {
		return nil, err
	}

	if allocation.UserID != userObjID {
		return nil, errors.New("unauthorized")
	}

	// Update pocket if provided
	if req.PocketID != "" {
		pocketObjID, err := primitive.ObjectIDFromHex(req.PocketID)
		if err != nil {
			return nil, errors.New("invalid pocket id")
		}

		pocket, err := s.pocketRepo.GetPocketByID(ctx, pocketObjID)
		if err != nil {
			return nil, errors.New("pocket not found")
		}

		if pocket.UserID != userObjID {
			return nil, errors.New("unauthorized: pocket does not belong to user")
		}

		allocation.PocketID = &pocketObjID
	}

	// Update user platform if provided
	if req.UserPlatformID != "" {
		userPlatformObjID, err := primitive.ObjectIDFromHex(req.UserPlatformID)
		if err != nil {
			return nil, errors.New("invalid user platform id")
		}

		userPlatform, err := s.userPlatformRepo.GetUserPlatformByID(ctx, userPlatformObjID)
		if err != nil {
			return nil, errors.New("user platform not found")
		}

		if userPlatform.UserID != userObjID {
			return nil, errors.New("unauthorized: user platform does not belong to user")
		}

		allocation.UserPlatformID = &userPlatformObjID
	}

	if req.Priority != nil {
		if !IsValidPriority(*req.Priority) {
			return nil, errors.New("invalid priority")
		}
		allocation.Priority = *req.Priority
	}

	if req.AllocationType != "" {
		if !IsValidAllocationType(req.AllocationType) {
			return nil, errors.New("invalid allocation type")
		}
		allocation.AllocationType = req.AllocationType
	}

	if req.Nominal != nil {
		if allocation.AllocationType == string(TypePercentage) && *req.Nominal > 100 {
			return nil, errors.New("percentage cannot exceed 100")
		}
		allocation.Nominal = *req.Nominal
	}

	if req.IsActive != nil {
		allocation.IsActive = *req.IsActive
	}

	err = s.repo.UpdateAllocation(ctx, allocationObjID, allocation)
	if err != nil {
		return nil, err
	}

	return allocation, nil
}

func (s *Service) DeleteAllocation(ctx context.Context, userID string, allocationID string) error {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user id")
	}

	allocationObjID, err := primitive.ObjectIDFromHex(allocationID)
	if err != nil {
		return errors.New("invalid allocation id")
	}

	allocation, err := s.repo.GetAllocationByID(ctx, allocationObjID)
	if err != nil {
		return err
	}

	if allocation.UserID != userObjID {
		return errors.New("unauthorized")
	}

	return s.repo.DeleteAllocation(ctx, allocationObjID)
}
