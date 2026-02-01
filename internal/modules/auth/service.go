package auth

import (
	"context"
	"errors"
	"log"
	"sort"

	"github.com/HasanNugroho/coin-be/internal/core/utils"
	authDTO "github.com/HasanNugroho/coin-be/internal/modules/auth/dto"
	"github.com/HasanNugroho/coin-be/internal/modules/category_template"
	"github.com/HasanNugroho/coin-be/internal/modules/pocket"
	"github.com/HasanNugroho/coin-be/internal/modules/pocket_template"
	"github.com/HasanNugroho/coin-be/internal/modules/user"
	"github.com/HasanNugroho/coin-be/internal/modules/user_category"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	userRepo             *user.Repository
	pocketRepo           *pocket.Repository
	pocketTemplateRepo   *pocket_template.Repository
	categoryTemplateRepo *category_template.Repository
	userCategoryRepo     *user_category.Repository
	redis                *redis.Client
	jwtManager           *utils.JWTManager
	passwordMgr          *utils.PasswordManager
}

func NewService(userRepo *user.Repository, pocketRepo *pocket.Repository, pocketTemplateRepo *pocket_template.Repository, categoryTemplateRepo *category_template.Repository, userCategoryRepo *user_category.Repository, redis *redis.Client, jwtManager *utils.JWTManager, passwordMgr *utils.PasswordManager) *Service {
	return &Service{
		userRepo:             userRepo,
		pocketRepo:           pocketRepo,
		pocketTemplateRepo:   pocketTemplateRepo,
		categoryTemplateRepo: categoryTemplateRepo,
		userCategoryRepo:     userCategoryRepo,
		redis:                redis,
		jwtManager:           jwtManager,
		passwordMgr:          passwordMgr,
	}
}

func (s *Service) Register(ctx context.Context, req *authDTO.RegisterRequest) (*user.User, error) {
	existingUser, _ := s.userRepo.GetUserByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	salt, err := s.passwordMgr.GenerateSalt()
	if err != nil {
		return nil, err
	}
	passwordHash, err := s.passwordMgr.HashPassword(req.Password, salt)
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
		PasswordHash: passwordHash,
		Salt:         salt,
		Name:         req.Name,
		Role:         userRole,
		IsActive:     true,
	}

	err = s.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, err
	}

	userProfile := &user.UserProfile{
		UserID:      newUser.ID,
		Phone:       req.Phone,
		TelegramId:  "",
		BaseSalary:  0,
		SalaryCycle: "monthly",
		SalaryDay:   1,
		PayCurrency: user.CurrencyIDR,
		Lang:        user.LanguageID,
		IsActive:    true,
	}

	err = s.userRepo.CreateUserProfile(ctx, userProfile)
	if err != nil {
		return nil, err
	}

	// Create default pockets from active templates
	err = s.createDefaultPockets(ctx, newUser.ID)
	if err != nil {
		// Rollback user creation on pocket creation failure
		_ = s.userRepo.DeleteUser(ctx, newUser.ID)
		return nil, err
	}

	// Create default category from active category templates

	err = s.createDefaultCategories(ctx, newUser.ID)
	if err != nil {
		// Rollback user creation on category creation failure
		_ = s.userRepo.DeleteUser(ctx, newUser.ID)
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

	if !s.passwordMgr.VerifyPassword(userRecord.PasswordHash, req.Password, userRecord.Salt) {
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

func (s *Service) createDefaultPockets(ctx context.Context, userID primitive.ObjectID) error {
	// Fetch all active pocket templates sorted by order
	templates, err := s.pocketTemplateRepo.GetActiveTemplatesSorted(ctx)
	if err != nil {
		log.Printf("failed to fetch active pocket templates for user %s: %v", userID.Hex(), err)
		return errors.New("failed to fetch pocket templates")
	}

	if len(templates) == 0 {
		log.Printf("no active pocket templates found for user %s", userID.Hex())
		return errors.New("no active pocket templates available")
	}

	// Validate that at least one MAIN pocket template exists
	hasMainTemplate := false
	var mainTemplate *pocket_template.PocketTemplate
	for _, t := range templates {
		if t.Type == string(pocket_template.TypeMain) {
			hasMainTemplate = true
			if mainTemplate == nil {
				mainTemplate = t
			}
		}
	}

	if !hasMainTemplate {
		log.Printf("no MAIN pocket template found for user %s", userID.Hex())
		return errors.New("no MAIN pocket template configured")
	}

	// If multiple MAIN templates, use the one with lowest order
	var mainTemplates []*pocket_template.PocketTemplate
	for _, t := range templates {
		if t.Type == string(pocket_template.TypeMain) && t.IsDefault {
			mainTemplates = append(mainTemplates, t)
		}
	}

	if len(mainTemplates) > 1 {
		sort.Slice(mainTemplates, func(i, j int) bool {
			return mainTemplates[i].Order < mainTemplates[j].Order
		})
		log.Printf("warning: multiple MAIN pocket templates found for user %s, using lowest order", userID.Hex())
		mainTemplate = mainTemplates[0]
	}

	// Create pockets from templates
	for _, template := range templates {
		newPocket := &pocket.Pocket{
			UserID:          userID,
			Name:            template.Name,
			Type:            template.Type,
			CategoryID:      template.CategoryID,
			Balance:         utils.NewDecimal128FromFloat(0),
			IsDefault:       template.IsDefault,
			IsActive:        true,
			IsLocked:        false,
			Icon:            template.Icon,
			IconColor:       template.IconColor,
			BackgroundColor: template.BackgroundColor,
		}

		err := s.pocketRepo.CreatePocket(ctx, newPocket)
		if err != nil {
			log.Printf("failed to create pocket from template %s for user %s: %v", template.ID.Hex(), userID.Hex(), err)
			return errors.New("failed to create default pockets")
		}
	}

	return nil
}

func (s *Service) createDefaultCategories(ctx context.Context, userID primitive.ObjectID) error {
	templates, err := s.categoryTemplateRepo.GetDefaults(ctx)
	if err != nil {
		log.Printf("failed to fetch active category templates for user %s: %v", userID.Hex(), err)
		return errors.New("failed to fetch category templates")
	}

	// Create a map to track template IDs for parent reference resolution
	templateMap := make(map[string]*category_template.CategoryTemplate)
	for _, template := range templates {
		templateMap[template.ID.Hex()] = template
	}

	// Create user categories from templates, handling parent relationships
	for _, template := range templates {
		// Convert TransactionType from category_template to user_category
		var transactionType *user_category.TransactionType
		if template.TransactionType != nil {
			tt := user_category.TransactionType(*template.TransactionType)
			transactionType = &tt
		}

		userCategory := &user_category.UserCategory{
			UserID:          userID,
			TemplateID:      &template.ID,
			Name:            template.Name,
			TransactionType: transactionType,
			Description:     template.Description,
			Icon:            template.Icon,
			Color:           template.Color,
			IsDefault:       template.IsDefault,
		}

		// Handle parent category reference
		if template.ParentID != nil {
			// Check if parent template exists in the defaults
			parentTemplate, exists := templateMap[template.ParentID.Hex()]
			if exists {
				// Find the corresponding user category for the parent template
				parentUserCategories, err := s.userCategoryRepo.FindAllByUserID(ctx, userID)
				if err == nil {
					for _, uc := range parentUserCategories {
						if uc.TemplateID != nil && uc.TemplateID.Hex() == parentTemplate.ID.Hex() {
							userCategory.ParentID = &uc.ID
							break
						}
					}
				}
			}
		}

		err := s.userCategoryRepo.Create(ctx, userCategory)
		if err != nil {
			log.Printf("failed to create user category from template %s for user %s: %v", template.ID.Hex(), userID.Hex(), err)
			return errors.New("failed to create default categories")
		}
	}

	return nil
}
