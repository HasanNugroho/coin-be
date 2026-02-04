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

func (s *Service) GetUserProfile(ctx context.Context, id string) (*dto.UserResponse, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	user, err := s.repo.GetUserByID(ctx, objID)
	if err != nil {
		return nil, err
	}

	profile, err := s.repo.GetUserProfileByUserID(ctx, objID)
	if err != nil {
		return nil, err
	}

	resp := &dto.UserResponse{
		ID:        user.ID.Hex(),
		Name:      user.Name,
		Email:     user.Email,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if profile != nil {
		resp.Phone = profile.Phone
		resp.TelegramId = profile.TelegramId
		resp.Currency = profile.PayCurrency
		resp.BaseSalary = profile.BaseSalary
		resp.SalaryCycle = profile.SalaryCycle
		resp.SalaryDay = profile.SalaryDay
		resp.Language = profile.Lang
		resp.AutoInputPayroll = profile.AutoInputPayroll
		if profile.DefaultUserPlatformID != nil {
			id := profile.DefaultUserPlatformID.Hex()
			resp.DefaultUserPlatformID = &id
		}
	}

	return resp, nil
}

func (s *Service) UpdateUser(ctx context.Context, id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
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
	if req.Email != "" {
		user.Email = req.Email
	}

	err = s.repo.UpdateUser(ctx, objID, user)
	if err != nil {
		return nil, err
	}

	profile, err := s.repo.GetUserProfileByUserID(ctx, objID)
	if err == nil && profile != nil {
		if req.Phone != "" {
			profile.Phone = req.Phone
		}
		if req.TelegramId != "" {
			profile.TelegramId = req.TelegramId
		}
		if req.Currency != "" {
			profile.PayCurrency = req.Currency
		}
		if req.BaseSalary > 0 {
			profile.BaseSalary = req.BaseSalary
		}
		if req.SalaryCycle != "" {
			profile.SalaryCycle = req.SalaryCycle
		}
		if req.SalaryDay > 0 {
			profile.SalaryDay = req.SalaryDay
		}
		if req.Language != "" {
			profile.Lang = req.Language
		}
		if req.AutoInputPayroll != nil {
			profile.AutoInputPayroll = *req.AutoInputPayroll
		}
		if req.DefaultUserPlatformID != "" {
			userPlatformID, err := primitive.ObjectIDFromHex(req.DefaultUserPlatformID)
			if err != nil {
				return nil, errors.New("invalid default_user_platform_id")
			}
			profile.DefaultUserPlatformID = &userPlatformID
		}

		err = s.repo.UpdateUserProfile(ctx, objID, profile)
		if err != nil {
			return nil, err
		}
	}

	resp := &dto.UserResponse{
		ID:        user.ID.Hex(),
		Name:      user.Name,
		Email:     user.Email,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if profile != nil {
		resp.Phone = profile.Phone
		resp.TelegramId = profile.TelegramId
		resp.Currency = profile.PayCurrency
		resp.BaseSalary = profile.BaseSalary
		resp.SalaryCycle = profile.SalaryCycle
		resp.SalaryDay = profile.SalaryDay
		resp.Language = profile.Lang
		resp.AutoInputPayroll = profile.AutoInputPayroll
		if profile.DefaultUserPlatformID != nil {
			id := profile.DefaultUserPlatformID.Hex()
			resp.DefaultUserPlatformID = &id
		}
	}

	return resp, nil
}

func (s *Service) DeleteUser(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid user id")
	}
	return s.repo.DeleteUser(ctx, objID)
}

func (s *Service) ListUsers(ctx context.Context, limit, skip int64, role, search, sort, order string) ([]*User, int64, error) {
	return s.repo.ListUsers(ctx, limit, skip, role, search, sort, order)
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
