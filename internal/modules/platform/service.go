package platform

import (
	"context"
	"errors"

	"github.com/HasanNugroho/coin-be/internal/modules/platform/dto"
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

func (s *Service) CreatePlatform(ctx context.Context, req *dto.CreatePlatformRequest) (*Platform, error) {
	if req.Name == "" {
		return nil, errors.New("platform name is required")
	}

	if !IsValidPlatformType(req.Type) {
		return nil, errors.New("invalid platform type")
	}

	existing, _ := s.repo.GetPlatformByName(ctx, req.Name)
	if existing != nil {
		return nil, errors.New("platform name already exists")
	}

	platform := &Platform{
		Name:     req.Name,
		Type:     req.Type,
		IsActive: req.IsActive,
	}

	err := s.repo.CreatePlatform(ctx, platform)
	if err != nil {
		return nil, err
	}

	return platform, nil
}

func (s *Service) GetPlatformByID(ctx context.Context, id string) (*Platform, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid platform id")
	}
	return s.repo.GetPlatformByID(ctx, objID)
}

func (s *Service) UpdatePlatform(ctx context.Context, id string, req *dto.UpdatePlatformRequest) (*Platform, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid platform id")
	}

	platform, err := s.repo.GetPlatformByID(ctx, objID)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		existing, _ := s.repo.GetPlatformByName(ctx, req.Name)
		if existing != nil && existing.ID != platform.ID {
			return nil, errors.New("platform name already exists")
		}
		platform.Name = req.Name
	}

	if req.Type != "" {
		if !IsValidPlatformType(req.Type) {
			return nil, errors.New("invalid platform type")
		}
		platform.Type = req.Type
	}

	if req.IsActive != nil {
		platform.IsActive = *req.IsActive
	}

	err = s.repo.UpdatePlatform(ctx, objID, platform)
	if err != nil {
		return nil, err
	}

	return platform, nil
}

func (s *Service) DeletePlatform(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid platform id")
	}
	return s.repo.DeletePlatform(ctx, objID)
}

func (s *Service) ListPlatforms(ctx context.Context, limit, skip int64) ([]*Platform, error) {
	return s.repo.ListPlatforms(ctx, limit, skip)
}

func (s *Service) ListActivePlatforms(ctx context.Context, limit, skip int64) ([]*Platform, error) {
	return s.repo.ListActivePlatforms(ctx, limit, skip)
}

func (s *Service) ListPlatformsByType(ctx context.Context, platformType string, limit, skip int64) ([]*Platform, error) {
	if !IsValidPlatformType(platformType) {
		return nil, errors.New("invalid platform type")
	}
	return s.repo.ListPlatformsByType(ctx, platformType, limit, skip)
}
