package category

import (
	"context"
	"errors"

	"github.com/HasanNugroho/coin-be/internal/modules/category/dto"
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

func (s *Service) CreateCategory(ctx context.Context, req *dto.CreateCategoryRequest) (*Category, error) {
	if req.Name == "" {
		return nil, errors.New("category name is required")
	}

	if req.Type != TypeTransaction && req.Type != TypePocket {
		return nil, errors.New("invalid category type")
	}

	existing, _ := s.repo.GetCategoryByName(ctx, req.Name)
	if existing != nil {
		return nil, errors.New("category name already exists")
	}

	var parentID *primitive.ObjectID
	if req.ParentID != "" {
		objID, err := primitive.ObjectIDFromHex(req.ParentID)
		if err != nil {
			return nil, errors.New("invalid parent id")
		}
		parent, err := s.repo.GetCategoryByID(ctx, objID)
		if err != nil {
			return nil, errors.New("parent category not found")
		}
		parentID = &parent.ID
	}

	category := &Category{
		Name:            req.Name,
		Type:            CategoryType(req.Type),
		TransactionType: (*TransactionType)(req.TransactionType),
		IsDefault:       req.IsDefault,
		ParentID:        parentID,
		Description:     req.Description,
		Icon:            req.Icon,
		Color:           req.Color,
	}

	err := s.repo.CreateCategory(ctx, category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *Service) GetCategoryByID(ctx context.Context, id string) (*Category, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid category id")
	}
	return s.repo.GetCategoryByID(ctx, objID)
}

func (s *Service) UpdateCategory(ctx context.Context, id string, req *dto.UpdateCategoryRequest) (*Category, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid category id")
	}

	category, err := s.repo.GetCategoryByID(ctx, objID)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		existing, _ := s.repo.GetCategoryByName(ctx, req.Name)
		if existing != nil && existing.ID != category.ID {
			return nil, errors.New("category name already exists")
		}
		category.Name = req.Name
	}

	if req.Type != "" {
		if req.Type != TypeTransaction && req.Type != TypePocket {
			return nil, errors.New("invalid category type")
		}
		category.Type = CategoryType(req.Type)
	}

	if req.TransactionType != nil {
		category.TransactionType = (*TransactionType)(req.TransactionType)
	}

	if req.Color != nil {
		category.Color = req.Color
	}

	if req.Icon != nil {
		category.Icon = req.Icon
	}

	if req.Description != nil {
		category.Description = req.Description
	}

	if req.IsDefault {
		category.IsDefault = req.IsDefault
	}

	if req.ParentID != "" {
		parentObjID, err := primitive.ObjectIDFromHex(req.ParentID)
		if err != nil {
			return nil, errors.New("invalid parent id")
		}
		parent, err := s.repo.GetCategoryByID(ctx, parentObjID)
		if err != nil {
			return nil, errors.New("parent category not found")
		}
		category.ParentID = &parent.ID
	}

	err = s.repo.UpdateCategory(ctx, objID, category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *Service) DeleteCategory(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid category id")
	}
	return s.repo.DeleteCategory(ctx, objID)
}

func (s *Service) ListCategories(ctx context.Context, limit, skip int64) ([]*Category, error) {
	return s.repo.ListCategories(ctx, limit, skip)
}

func (s *Service) ListCategoriesByType(ctx context.Context, categoryType string, limit, skip int64) ([]*Category, error) {
	if categoryType != TypeTransaction && categoryType != TypePocket {
		return nil, errors.New("invalid category type")
	}
	return s.repo.ListCategoriesByType(ctx, categoryType, limit, skip)
}

func (s *Service) ListSubcategories(ctx context.Context, parentID string) ([]*Category, error) {
	objID, err := primitive.ObjectIDFromHex(parentID)
	if err != nil {
		return nil, errors.New("invalid parent id")
	}
	return s.repo.ListSubcategories(ctx, objID)
}
