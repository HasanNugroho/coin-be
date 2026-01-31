package pocket_template

import (
	"context"
	"errors"

	"github.com/HasanNugroho/coin-be/internal/modules/pocket_template/dto"
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

func (s *Service) CreatePocketTemplate(ctx context.Context, req *dto.CreatePocketTemplateRequest) (*PocketTemplate, error) {
	categoryID, err := primitive.ObjectIDFromHex(req.CategoryID)
	if err != nil {
		return nil, errors.New("invalid category id")
	}

	existing, _ := s.repo.GetPocketTemplateByName(ctx, req.Name)
	if existing != nil {
		return nil, errors.New("pocket template name already exists")
	}

	template := &PocketTemplate{
		Name:       req.Name,
		Type:       req.Type,
		CategoryID: &categoryID,
		IsDefault:  req.IsDefault,
		IsActive:   req.IsActive,
		Order:      req.Order,
	}

	err = s.repo.CreatePocketTemplate(ctx, template)
	if err != nil {
		return nil, err
	}

	return template, nil
}

func (s *Service) GetPocketTemplateByID(ctx context.Context, id string) (*PocketTemplate, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid pocket template id")
	}
	return s.repo.GetPocketTemplateByID(ctx, objID)
}

func (s *Service) UpdatePocketTemplate(ctx context.Context, id string, req *dto.UpdatePocketTemplateRequest) (*PocketTemplate, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid pocket template id")
	}

	template, err := s.repo.GetPocketTemplateByID(ctx, objID)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		existing, _ := s.repo.GetPocketTemplateByName(ctx, req.Name)
		if existing != nil && existing.ID != template.ID {
			return nil, errors.New("pocket template name already exists")
		}
		template.Name = req.Name
	}

	if req.Type != "" {
		template.Type = req.Type
	}

	if req.CategoryID != "" {
		categoryID, err := primitive.ObjectIDFromHex(req.CategoryID)
		if err != nil {
			return nil, errors.New("invalid category id")
		}
		template.CategoryID = &categoryID
	}

	if req.Icon != "" {
		template.Icon = &req.Icon
	}

	if req.IconColor != "" {
		template.IconColor = &req.IconColor
	}

	if req.BackgroundColor != "" {
		template.BackgroundColor = &req.BackgroundColor
	}

	template.IsDefault = req.IsDefault
	template.IsActive = req.IsActive
	template.Order = req.Order

	err = s.repo.UpdatePocketTemplate(ctx, objID, template)
	if err != nil {
		return nil, err
	}

	return template, nil
}

func (s *Service) DeletePocketTemplate(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid pocket template id")
	}
	return s.repo.DeletePocketTemplate(ctx, objID)
}

func (s *Service) ListPocketTemplates(ctx context.Context, limit, skip int64) ([]*PocketTemplate, error) {
	return s.repo.ListPocketTemplates(ctx, limit, skip)
}

func (s *Service) ListActivePocketTemplates(ctx context.Context, limit, skip int64) ([]*PocketTemplate, error) {
	return s.repo.ListActivePocketTemplates(ctx, limit, skip)
}

func (s *Service) ListPocketTemplatesByType(ctx context.Context, templateType string, limit, skip int64) ([]*PocketTemplate, error) {
	return s.repo.ListPocketTemplatesByType(ctx, templateType, limit, skip)
}
