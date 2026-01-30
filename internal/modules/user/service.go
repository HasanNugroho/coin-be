package user

import (
	"context"
	"errors"

	"github.com/HasanNugroho/coin-be/internal/modules/user/dto"
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

func (s *Service) GetUserByID(ctx context.Context, id string) (*User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user id")
	}
	return s.repo.GetUserByID(ctx, objID)
}

func (s *Service) UpdateUser(ctx context.Context, id string, req *dto.UpdateUserRequest) (*User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	user, err := s.repo.GetUserByID(ctx, objID)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	err = s.repo.UpdateUser(ctx, objID, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) DeleteUser(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid user id")
	}
	return s.repo.DeleteUser(ctx, objID)
}

func (s *Service) ListUsers(ctx context.Context, limit, skip int64) ([]*User, error) {
	return s.repo.ListUsers(ctx, limit, skip)
}

func (s *Service) CreateUserProfile(ctx context.Context, userID string, req *dto.CreateUserProfileRequest) (*UserProfile, error) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	user, err := s.repo.GetUserByID(ctx, objID)
	if err != nil {
		return nil, err
	}

	profile := &UserProfile{
		UserID:      user.ID,
		BaseSalary:  req.BaseSalary,
		SalaryCycle: req.SalaryCycle,
		SalaryDay:   req.SalaryDay,
		PayCurrency: req.PayCurrency,
		IsActive:    true,
	}

	err = s.repo.CreateUserProfile(ctx, profile)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *Service) DisableUser(ctx context.Context, userID string) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user id")
	}

	user, err := s.repo.GetUserByID(ctx, objID)
	if err != nil {
		return err
	}

	user.IsActive = false
	return s.repo.UpdateUser(ctx, objID, user)
}

func (s *Service) EnableUser(ctx context.Context, userID string) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user id")
	}

	user, err := s.repo.GetUserByID(ctx, objID)
	if err != nil {
		return err
	}

	user.IsActive = true
	return s.repo.UpdateUser(ctx, objID, user)
}
