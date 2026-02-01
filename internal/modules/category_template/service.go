package category_template

import (
	"context"
	"errors"

	"github.com/HasanNugroho/coin-be/internal/modules/category_template/dto"
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

func (s *Service) CreateCategoryTemplate(ctx context.Context, req *dto.CreateCategoryTemplateRequest) (*CategoryTemplate, error) {
	if req.Name == "" {
		return nil, errors.New("category template name is required")
	}

	var parentID *primitive.ObjectID
	if req.ParentID != "" {
		objID, err := primitive.ObjectIDFromHex(req.ParentID)
		if err != nil {
			return nil, errors.New("invalid parent id")
		}
		parent, err := s.repo.FindByID(ctx, objID)
		if err != nil {
			return nil, errors.New("parent category template not found")
		}
		parentID = &parent.ID
	}

	template := &CategoryTemplate{
		Name:            req.Name,
		TransactionType: (*TransactionType)(req.TransactionType),
		IsDefault:       req.IsDefault,
		ParentID:        parentID,
		Description:     req.Description,
		Icon:            req.Icon,
		Color:           req.Color,
	}

	err := s.repo.Create(ctx, template)
	if err != nil {
		return nil, err
	}

	return template, nil
}

func (s *Service) GetCategoryTemplateByID(ctx context.Context, id string) (*CategoryTemplate, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid category template id")
	}
	return s.repo.FindByID(ctx, objID)
}

func (s *Service) FindAll(ctx context.Context) ([]*CategoryTemplate, error) {
	templates, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return templates, nil
}

func (s *Service) FindAllWithFilter(ctx context.Context, transactionType *string, page int64, pageSize int64) ([]*CategoryTemplate, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	templates, total, err := s.repo.FindAllWithFilter(ctx, transactionType, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return templates, total, nil
}

func (s *Service) FindAllParent(ctx context.Context, transactionType *string) ([]*CategoryTemplate, error) {
	templates, err := s.repo.FindAllParent(ctx, transactionType)
	if err != nil {
		return nil, err
	}
	return templates, nil
}

func (s *Service) UpdateCategoryTemplate(ctx context.Context, id string, req *dto.UpdateCategoryTemplateRequest) (*CategoryTemplate, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid category template id")
	}

	template, err := s.repo.FindByID(ctx, objID)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		template.Name = req.Name
	}

	if req.TransactionType != nil {
		template.TransactionType = (*TransactionType)(req.TransactionType)
	}

	if req.Color != nil {
		template.Color = req.Color
	}

	if req.Icon != nil {
		template.Icon = req.Icon
	}

	if req.Description != nil {
		template.Description = req.Description
	}

	if req.IsDefault {
		template.IsDefault = req.IsDefault
	}

	if req.ParentID != "" {
		parentObjID, err := primitive.ObjectIDFromHex(req.ParentID)
		if err != nil {
			return nil, errors.New("invalid parent id")
		}
		parent, err := s.repo.FindByID(ctx, parentObjID)
		if err != nil {
			return nil, errors.New("parent category template not found")
		}
		template.ParentID = &parent.ID
	}

	err = s.repo.Update(ctx, objID, template)
	if err != nil {
		return nil, err
	}

	return template, nil
}

func (s *Service) DeleteCategoryTemplate(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid category template id")
	}
	return s.repo.SoftDelete(ctx, objID)
}
