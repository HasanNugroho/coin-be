package user_category

import (
	"context"
	"errors"

	"github.com/HasanNugroho/coin-be/internal/modules/user_category/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	repo                  *Repository
	categoryTemplateDB    *mongo.Database
}

func NewService(r *Repository, db *mongo.Database) *Service {
	return &Service{
		repo:               r,
		categoryTemplateDB: db,
	}
}

func (s *Service) CreateUserCategory(ctx context.Context, userID string, req *dto.CreateUserCategoryRequest) (*UserCategory, error) {
	if req.Name == "" {
		return nil, errors.New("user category name is required")
	}

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	var templateID *primitive.ObjectID
	if req.TemplateID != "" {
		objID, err := primitive.ObjectIDFromHex(req.TemplateID)
		if err != nil {
			return nil, errors.New("invalid template id")
		}
		templateCollection := s.categoryTemplateDB.Collection("category_templates")
		count, err := templateCollection.CountDocuments(ctx, primitive.M{"_id": objID, "is_deleted": false})
		if err != nil {
			return nil, err
		}
		if count == 0 {
			return nil, errors.New("category template not found")
		}
		templateID = &objID
	}

	var parentID *primitive.ObjectID
	if req.ParentID != "" {
		objID, err := primitive.ObjectIDFromHex(req.ParentID)
		if err != nil {
			return nil, errors.New("invalid parent id")
		}
		parent, err := s.repo.FindByID(ctx, objID, userObjID)
		if err != nil {
			return nil, errors.New("parent user category not found")
		}
		parentID = &parent.ID
	}

	category := &UserCategory{
		UserID:          userObjID,
		TemplateID:      templateID,
		Name:            req.Name,
		TransactionType: (*TransactionType)(req.TransactionType),
		IsDefault:       req.IsDefault,
		ParentID:        parentID,
		Description:     req.Description,
		Icon:            req.Icon,
		Color:           req.Color,
	}

	err = s.repo.Create(ctx, category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *Service) GetUserCategoryByID(ctx context.Context, id string, userID string) (*UserCategory, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user category id")
	}

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	return s.repo.FindByID(ctx, objID, userObjID)
}

func (s *Service) GetUserCategories(ctx context.Context, userID string) ([]*UserCategory, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	return s.repo.FindAllByUserID(ctx, userObjID)
}

func (s *Service) UpdateUserCategory(ctx context.Context, id string, userID string, req *dto.UpdateUserCategoryRequest) (*UserCategory, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user category id")
	}

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	category, err := s.repo.FindByID(ctx, objID, userObjID)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		category.Name = req.Name
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

	if req.TemplateID != "" {
		templateObjID, err := primitive.ObjectIDFromHex(req.TemplateID)
		if err != nil {
			return nil, errors.New("invalid template id")
		}
		templateCollection := s.categoryTemplateDB.Collection("category_templates")
		count, err := templateCollection.CountDocuments(ctx, primitive.M{"_id": templateObjID, "is_deleted": false})
		if err != nil {
			return nil, err
		}
		if count == 0 {
			return nil, errors.New("category template not found")
		}
		category.TemplateID = &templateObjID
	}

	if req.ParentID != "" {
		parentObjID, err := primitive.ObjectIDFromHex(req.ParentID)
		if err != nil {
			return nil, errors.New("invalid parent id")
		}
		parent, err := s.repo.FindByID(ctx, parentObjID, userObjID)
		if err != nil {
			return nil, errors.New("parent user category not found")
		}
		category.ParentID = &parent.ID
	}

	err = s.repo.Update(ctx, objID, userObjID, category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *Service) DeleteUserCategory(ctx context.Context, id string, userID string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid user category id")
	}

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user id")
	}

	return s.repo.SoftDelete(ctx, objID, userObjID)
}
