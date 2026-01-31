package pocket

import (
	"context"
	"errors"

	"github.com/HasanNugroho/coin-be/internal/modules/pocket/dto"
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

func (s *Service) CreatePocket(ctx context.Context, userID string, req *dto.CreatePocketRequest) (*Pocket, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	if req.Type == string(TypeMain) {
		existing, _ := s.repo.GetMainPocketByUserID(ctx, userObjID)
		if existing != nil {
			return nil, errors.New("user already has a main pocket")
		}
	}

	var categoryID *primitive.ObjectID
	if req.CategoryID != "" {
		catID, err := primitive.ObjectIDFromHex(req.CategoryID)
		if err != nil {
			return nil, errors.New("invalid category id")
		}
		categoryID = &catID
	}

	pocket := &Pocket{
		UserID:          userObjID,
		Name:            req.Name,
		Type:            req.Type,
		CategoryID:      categoryID,
		Balance:         NewDecimal128(0),
		IsDefault:       req.Type == string(TypeMain),
		IsActive:        true,
		IsLocked:        false,
		Icon:            stringPtr(req.Icon),
		IconColor:       stringPtr(req.IconColor),
		BackgroundColor: stringPtr(req.BackgroundColor),
	}

	err = s.repo.CreatePocket(ctx, pocket)
	if err != nil {
		return nil, err
	}

	return pocket, nil
}

func (s *Service) GetPocketByID(ctx context.Context, userID string, pocketID string) (*Pocket, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	pocketObjID, err := primitive.ObjectIDFromHex(pocketID)
	if err != nil {
		return nil, errors.New("invalid pocket id")
	}

	pocket, err := s.repo.GetPocketByID(ctx, pocketObjID)
	if err != nil {
		return nil, err
	}

	if pocket.UserID != userObjID {
		return nil, errors.New("unauthorized")
	}

	return pocket, nil
}

func (s *Service) UpdatePocket(ctx context.Context, userID string, pocketID string, req *dto.UpdatePocketRequest) (*Pocket, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	pocketObjID, err := primitive.ObjectIDFromHex(pocketID)
	if err != nil {
		return nil, errors.New("invalid pocket id")
	}

	pocket, err := s.repo.GetPocketByID(ctx, pocketObjID)
	if err != nil {
		return nil, err
	}

	if pocket.UserID != userObjID {
		return nil, errors.New("unauthorized")
	}

	if pocket.IsLocked {
		return nil, errors.New("pocket is locked")
	}

	if pocket.Type == string(TypeMain) {
		return nil, errors.New("cannot update main pocket")
	}

	if req.Name != "" {
		pocket.Name = req.Name
	}

	if req.Type != "" && req.Type != pocket.Type {
		pocket.Type = req.Type
	}

	if req.CategoryID != "" {
		catID, err := primitive.ObjectIDFromHex(req.CategoryID)
		if err != nil {
			return nil, errors.New("invalid category id")
		}
		pocket.CategoryID = &catID
	}

	if req.Icon != "" {
		pocket.Icon = stringPtr(req.Icon)
	}

	if req.IconColor != "" {
		pocket.IconColor = stringPtr(req.IconColor)
	}

	if req.BackgroundColor != "" {
		pocket.BackgroundColor = stringPtr(req.BackgroundColor)
	}

	if req.IsActive != nil {
		pocket.IsActive = *req.IsActive
	}

	err = s.repo.UpdatePocket(ctx, pocketObjID, pocket)
	if err != nil {
		return nil, err
	}

	return pocket, nil
}

func (s *Service) DeletePocket(ctx context.Context, userID string, pocketID string) error {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user id")
	}

	pocketObjID, err := primitive.ObjectIDFromHex(pocketID)
	if err != nil {
		return errors.New("invalid pocket id")
	}

	pocket, err := s.repo.GetPocketByID(ctx, pocketObjID)
	if err != nil {
		return err
	}

	if pocket.UserID != userObjID {
		return errors.New("unauthorized")
	}

	if pocket.Type == string(TypeMain) {
		return errors.New("cannot delete main pocket")
	}

	if pocket.IsLocked {
		return errors.New("pocket is locked")
	}

	return s.repo.DeletePocket(ctx, pocketObjID)
}

func (s *Service) GetUserPockets(ctx context.Context, userID string, limit int64, skip int64) ([]*Pocket, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	return s.repo.GetPocketsByUserID(ctx, userObjID, limit, skip)
}

func (s *Service) GetActiveUserPockets(ctx context.Context, userID string, limit int64, skip int64) ([]*Pocket, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	return s.repo.GetActivePocketsByUserID(ctx, userObjID, limit, skip)
}

func (s *Service) GetMainPocket(ctx context.Context, userID string) (*Pocket, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	pocket, err := s.repo.GetMainPocketByUserID(ctx, userObjID)
	if err != nil {
		return nil, err
	}

	if pocket == nil {
		return nil, errors.New("main pocket not found")
	}

	return pocket, nil
}

func (s *Service) CreateSystemPocket(ctx context.Context, userID string, req *dto.CreateSystemPocketRequest) (*Pocket, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	var categoryID *primitive.ObjectID
	if req.CategoryID != "" {
		catID, err := primitive.ObjectIDFromHex(req.CategoryID)
		if err != nil {
			return nil, errors.New("invalid category id")
		}
		categoryID = &catID
	}

	pocket := &Pocket{
		UserID:          userObjID,
		Name:            req.Name,
		Type:            string(TypeSystem),
		CategoryID:      categoryID,
		Balance:         NewDecimal128(0),
		IsDefault:       false,
		IsActive:        true,
		IsLocked:        true,
		Icon:            stringPtr(req.Icon),
		IconColor:       stringPtr(req.IconColor),
		BackgroundColor: stringPtr(req.BackgroundColor),
	}

	err = s.repo.CreatePocket(ctx, pocket)
	if err != nil {
		return nil, err
	}

	return pocket, nil
}

func (s *Service) GetAllPockets(ctx context.Context, limit int64, skip int64) ([]*Pocket, error) {
	return s.repo.GetAllPockets(ctx, limit, skip)
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
