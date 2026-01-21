package auth

import (
	"context"
	"errors"

	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/internal/modules/allocation"
	authDTO "github.com/HasanNugroho/coin-be/internal/modules/auth/dto"
	"github.com/HasanNugroho/coin-be/internal/modules/category"
	"github.com/HasanNugroho/coin-be/internal/modules/user"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	userRepo       *user.Repository
	categoryRepo   *category.Repository
	allocationRepo *allocation.Repository
	redis          *redis.Client
	jwtManager     *utils.JWTManager
	passwordMgr    *utils.PasswordManager
}

func NewService(userRepo *user.Repository, categoryRepo *category.Repository, allocationRepo *allocation.Repository, redis *redis.Client, jwtManager *utils.JWTManager, passwordMgr *utils.PasswordManager) *Service {
	return &Service{
		userRepo:       userRepo,
		categoryRepo:   categoryRepo,
		allocationRepo: allocationRepo,
		redis:          redis,
		jwtManager:     jwtManager,
		passwordMgr:    passwordMgr,
	}
}

func (s *Service) Register(ctx context.Context, req *authDTO.RegisterRequest) (*user.User, error) {
	existingUser, _ := s.userRepo.GetUserByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	passwordHash, err := s.passwordMgr.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Assign role: first user is admin, subsequent users are regular users
	userRole := user.RoleUser
	adminCount, err := s.userRepo.CountUsersByRole(ctx, user.RoleAdmin)
	if err == nil && adminCount == 0 {
		userRole = user.RoleAdmin
	}

	newUser := &user.User{
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: passwordHash,
		Name:         req.Name,
		Role:         userRole,
		IsActive:     true,
	}

	err = s.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, err
	}

	if err := s.categoryRepo.CreateDefaultCategories(ctx, newUser.ID); err != nil {
		return nil, err
	}

	if err := s.allocationRepo.CreateDefaultAllocations(ctx, newUser.ID); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *Service) Login(ctx context.Context, req *authDTO.LoginRequest) (*authDTO.LoginResponse, error) {
	userRecord, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !userRecord.IsActive {
		return nil, errors.New("user account is inactive")
	}

	if !s.passwordMgr.VerifyPassword(userRecord.PasswordHash, req.Password) {
		return nil, errors.New("invalid email or password")
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(userRecord.ID.Hex(), userRecord.Email)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(userRecord.ID.Hex())
	if err != nil {
		return nil, err
	}

	err = s.storeRefreshToken(ctx, userRecord.ID.Hex(), refreshToken)
	if err != nil {
		return nil, err
	}

	return &authDTO.LoginResponse{
		User: userRecord,
		TokenPair: authDTO.TokenPair{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (s *Service) RefreshAccessToken(ctx context.Context, userID, refreshToken string) (string, error) {
	verifiedUserID, err := s.jwtManager.VerifyRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	if verifiedUserID != userID {
		return "", errors.New("refresh token mismatch")
	}

	storedToken, err := s.redis.Get(ctx, "refresh_token:"+userID).Result()
	if err != nil {
		return "", errors.New("refresh token not found or expired")
	}

	if storedToken != refreshToken {
		return "", errors.New("refresh token mismatch")
	}

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return "", errors.New("invalid user id")
	}

	userRecord, err := s.userRepo.GetUserByID(ctx, objID)
	if err != nil {
		return "", err
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(userRecord.ID.Hex(), userRecord.Email)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s *Service) Logout(ctx context.Context, userID string) error {
	return s.redis.Del(ctx, "refresh_token:"+userID).Err()
}

func (s *Service) storeRefreshToken(ctx context.Context, userID, token string) error {
	return s.redis.Set(ctx, "refresh_token:"+userID, token, 0).Err()
}
