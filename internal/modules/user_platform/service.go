package user_platform

import (
	"context"
	"errors"

	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/internal/modules/platform"
	"github.com/HasanNugroho/coin-be/internal/modules/user_platform/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo         *UserPlatformRepository
	platformRepo *platform.Repository
}

func NewService(r *UserPlatformRepository, pr *platform.Repository) *Service {
	return &Service{
		repo:         r,
		platformRepo: pr,
	}
}

func (s *Service) CreateUserPlatform(ctx context.Context, userID string, req *dto.CreateUserPlatformRequest) (*UserPlatform, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	platformObjID, err := primitive.ObjectIDFromHex(req.PlatformID)
	if err != nil {
		return nil, errors.New("invalid platform id")
	}

	// Verify platform exists
	platform, err := s.platformRepo.GetPlatformByID(ctx, platformObjID)
	if err != nil {
		return nil, errors.New("platform not found")
	}

	if !platform.IsActive {
		return nil, errors.New("platform is not active")
	}

	userPlatform := &UserPlatform{
		UserID:     userObjID,
		PlatformID: platformObjID,
		AliasName:  req.AliasName,
		Balance:    utils.NewDecimal128FromFloat(0),
		IsActive:   true,
	}

	err = s.repo.CreateUserPlatform(ctx, userPlatform)
	if err != nil {
		return nil, err
	}

	return userPlatform, nil
}

func (s *Service) GetUserPlatformByID(ctx context.Context, userID string, userPlatformID string) (*UserPlatform, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	userPlatformObjID, err := primitive.ObjectIDFromHex(userPlatformID)
	if err != nil {
		return nil, errors.New("invalid user platform id")
	}

	userPlatform, err := s.repo.GetUserPlatformByID(ctx, userPlatformObjID)
	if err != nil {
		return nil, err
	}

	if userPlatform.UserID != userObjID {
		return nil, errors.New("unauthorized")
	}

	return userPlatform, nil
}

func (s *Service) ListUserPlatforms(ctx context.Context, userID string) ([]*UserPlatform, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	return s.repo.GetUserPlatformsByUserID(ctx, userObjID)
}

func (s *Service) ListUserPlatformsDropdown(ctx context.Context, userID string) ([]*UserPlatform, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	return s.repo.GetUserPlatformsByUserIDDropdown(ctx, userObjID)
}

func (s *Service) UpdateUserPlatform(ctx context.Context, userID string, userPlatformID string, req *dto.UpdateUserPlatformRequest) (*UserPlatform, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	userPlatformObjID, err := primitive.ObjectIDFromHex(userPlatformID)
	if err != nil {
		return nil, errors.New("invalid user platform id")
	}

	userPlatform, err := s.repo.GetUserPlatformByID(ctx, userPlatformObjID)
	if err != nil {
		return nil, err
	}

	if userPlatform.UserID != userObjID {
		return nil, errors.New("unauthorized")
	}

	if req.AliasName != nil {
		userPlatform.AliasName = req.AliasName
	}

	if req.IsActive != nil {
		userPlatform.IsActive = *req.IsActive
	}

	err = s.repo.UpdateUserPlatform(ctx, userPlatformObjID, userPlatform)
	if err != nil {
		return nil, err
	}

	return userPlatform, nil
}

func (s *Service) DeleteUserPlatform(ctx context.Context, userID string, userPlatformID string) error {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user id")
	}

	userPlatformObjID, err := primitive.ObjectIDFromHex(userPlatformID)
	if err != nil {
		return errors.New("invalid user platform id")
	}

	userPlatform, err := s.repo.GetUserPlatformByID(ctx, userPlatformObjID)
	if err != nil {
		return err
	}

	if userPlatform.UserID != userObjID {
		return errors.New("unauthorized")
	}

	return s.repo.DeleteUserPlatform(ctx, userPlatformObjID)
}
