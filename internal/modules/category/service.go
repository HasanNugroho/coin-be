package category

import (
	"context"
	"fmt"

	"github.com/HasanNugroho/coin-be/internal/modules/category/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateCategory(ctx context.Context, userID primitive.ObjectID, req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	category := &Category{
		UserID:    userID,
		Name:      req.Name,
		Type:      CategoryType(req.Type),
		Icon:      req.Icon,
		Color:     req.Color,
		IsDefault: false,
	}

	if err := s.repo.Create(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return &dto.CategoryResponse{
		ID:        category.ID,
		UserID:    category.UserID,
		Name:      category.Name,
		Type:      string(category.Type),
		Icon:      category.Icon,
		Color:     category.Color,
		IsDefault: category.IsDefault,
		CreatedAt: category.CreatedAt,
	}, nil
}

func (s *Service) GetCategories(ctx context.Context, userID primitive.ObjectID) ([]*dto.CategoryResponse, error) {
	categories, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	responses := make([]*dto.CategoryResponse, len(categories))
	for i, cat := range categories {
		responses[i] = &dto.CategoryResponse{
			ID:        cat.ID,
			UserID:    cat.UserID,
			Name:      cat.Name,
			Type:      string(cat.Type),
			Icon:      cat.Icon,
			Color:     cat.Color,
			IsDefault: cat.IsDefault,
			CreatedAt: cat.CreatedAt,
		}
	}

	return responses, nil
}

func (s *Service) GetCategoriesByType(ctx context.Context, userID primitive.ObjectID, categoryType string) ([]*dto.CategoryResponse, error) {
	categories, err := s.repo.GetByUserIDAndType(ctx, userID, CategoryType(categoryType))
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	responses := make([]*dto.CategoryResponse, len(categories))
	for i, cat := range categories {
		responses[i] = &dto.CategoryResponse{
			ID:        cat.ID,
			UserID:    cat.UserID,
			Name:      cat.Name,
			Type:      string(cat.Type),
			Icon:      cat.Icon,
			Color:     cat.Color,
			IsDefault: cat.IsDefault,
			CreatedAt: cat.CreatedAt,
		}
	}

	return responses, nil
}

func (s *Service) GetCategoryByID(ctx context.Context, id primitive.ObjectID) (*dto.CategoryResponse, error) {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return &dto.CategoryResponse{
		ID:        category.ID,
		UserID:    category.UserID,
		Name:      category.Name,
		Type:      string(category.Type),
		Icon:      category.Icon,
		Color:     category.Color,
		IsDefault: category.IsDefault,
		CreatedAt: category.CreatedAt,
	}, nil
}

func (s *Service) UpdateCategory(ctx context.Context, id primitive.ObjectID, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Icon != "" {
		category.Icon = req.Icon
	}
	if req.Color != "" {
		category.Color = req.Color
	}

	if err := s.repo.Update(ctx, id, category); err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return &dto.CategoryResponse{
		ID:        category.ID,
		UserID:    category.UserID,
		Name:      category.Name,
		Type:      string(category.Type),
		Icon:      category.Icon,
		Color:     category.Color,
		IsDefault: category.IsDefault,
		CreatedAt: category.CreatedAt,
	}, nil
}

func (s *Service) DeleteCategory(ctx context.Context, id primitive.ObjectID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}
	return nil
}
