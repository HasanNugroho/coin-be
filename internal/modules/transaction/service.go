package transaction

import (
	"context"
	"errors"
	"time"

	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/internal/modules/pocket"
	"github.com/HasanNugroho/coin-be/internal/modules/transaction/dto"
	"github.com/HasanNugroho/coin-be/internal/modules/user_platform"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo             *Repository
	pocketRepo       *pocket.Repository
	userPlatformRepo *user_platform.UserPlatformRepository
	balanceProcessor *BalanceProcessor
}

func NewService(r *Repository, pr *pocket.Repository, upr *user_platform.UserPlatformRepository) *Service {
	return &Service{
		repo:             r,
		pocketRepo:       pr,
		userPlatformRepo: upr,
		balanceProcessor: NewBalanceProcessor(pr, upr),
	}
}

func (s *Service) CreateTransaction(ctx context.Context, userID string, req *dto.CreateTransactionRequest) (*Transaction, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	if !IsValidTransactionType(req.Type) {
		return nil, errors.New("invalid transaction type")
	}

	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		return nil, errors.New("invalid date format")
	}

	var pocketFrom *primitive.ObjectID
	var pocketTo *primitive.ObjectID
	var userPlatformFrom *primitive.ObjectID
	var userPlatformTo *primitive.ObjectID

	if req.PocketFromID != "" {
		pocketFromID, err := primitive.ObjectIDFromHex(req.PocketFromID)
		if err != nil {
			return nil, errors.New("invalid pocket_from id")
		}
		pocketFrom = &pocketFromID
	}

	if req.PocketToID != "" {
		pocketToID, err := primitive.ObjectIDFromHex(req.PocketToID)
		if err != nil {
			return nil, errors.New("invalid pocket_to id")
		}
		pocketTo = &pocketToID
	}

	if req.UserPlatformFromID != "" {
		userPlatformFromID, err := primitive.ObjectIDFromHex(req.UserPlatformFromID)
		if err != nil {
			return nil, errors.New("invalid user_platform_from id")
		}
		userPlatformFrom = &userPlatformFromID
	}

	if req.UserPlatformToID != "" {
		userPlatformToID, err := primitive.ObjectIDFromHex(req.UserPlatformToID)
		if err != nil {
			return nil, errors.New("invalid user_platform_to id")
		}
		userPlatformTo = &userPlatformToID
	}

	var categoryID *primitive.ObjectID
	if req.CategoryID != "" {
		catID, err := primitive.ObjectIDFromHex(req.CategoryID)
		if err != nil {
			return nil, errors.New("invalid category id")
		}
		categoryID = &catID
	}

	// Validate transaction rules based on type and provided fields
	if err := s.validateTransactionRules(ctx, req.Type, userObjID, pocketFrom, pocketTo, userPlatformFrom, userPlatformTo); err != nil {
		return nil, err
	}

	// Validate ownership of all pockets
	if err := s.validatePocket(ctx, userObjID, pocketFrom, pocketTo, req.Amount); err != nil {
		return nil, err
	}

	// Validate ownership of all user platforms
	if err := s.validateUserPlatform(ctx, userObjID, userPlatformFrom, userPlatformTo, req.Amount); err != nil {
		return nil, err
	}

	transaction := &Transaction{
		UserID:             userObjID,
		Type:               req.Type,
		Amount:             req.Amount,
		PocketFromID:       pocketFrom,
		PocketToID:         pocketTo,
		UserPlatformFromID: userPlatformFrom,
		UserPlatformToID:   userPlatformTo,
		CategoryID:         categoryID,
		Note:               stringPtr(req.Note),
		Date:               date,
		Ref:                stringPtr(req.Ref),
	}

	// Create transaction record
	if err := s.repo.CreateTransaction(ctx, transaction); err != nil {
		return nil, err
	}

	// Process balance updates through centralized processor
	if err := s.balanceProcessor.ProcessTransaction(ctx, req.Type, req.Amount, pocketFrom, pocketTo, userPlatformFrom, userPlatformTo); err != nil {
		s.repo.DeleteTransaction(ctx, transaction.ID)
		return nil, err
	}

	return transaction, nil
}

func (s *Service) validateTransactionRules(ctx context.Context, txType string, userID primitive.ObjectID, pocketFrom, pocketTo, userPlatformFrom, userPlatformTo *primitive.ObjectID) error {
	switch txType {
	case string(TypeIncome):
		// Income must have destination (pocket or platform)
		if pocketTo == nil && userPlatformTo == nil {
			return errors.New("pocket_to or user_platform_to is required for INCOME transactions")
		}
		// Income cannot have source
		if pocketFrom != nil || userPlatformFrom != nil {
			return errors.New("pocket_from and user_platform_from must be null for INCOME transactions")
		}

	case string(TypeExpense):
		// Expense must have source (pocket or platform)
		if pocketFrom == nil && userPlatformFrom == nil {
			return errors.New("pocket_from or user_platform_from is required for EXPENSE transactions")
		}
		// Expense cannot have destination
		if pocketTo != nil || userPlatformTo != nil {
			return errors.New("pocket_to and user_platform_to must be null for EXPENSE transactions")
		}

	case string(TypeTransfer):
		// Transfer requires either (pocket-to-pocket) or (platform-to-platform) or (both pairs)
		hasPocketPair := pocketFrom != nil && pocketTo != nil
		hasPlatformPair := userPlatformFrom != nil && userPlatformTo != nil

		if !hasPocketPair && !hasPlatformPair {
			return errors.New("TRANSFER requires either (pocket_from + pocket_to) or (user_platform_from + user_platform_to) or both pairs")
		}

		// Validate pocket pair if present
		if hasPocketPair && pocketFrom.Hex() == pocketTo.Hex() {
			return errors.New("pocket_from and pocket_to cannot be the same")
		}

		// Validate platform pair if present
		if hasPlatformPair && userPlatformFrom.Hex() == userPlatformTo.Hex() {
			return errors.New("user_platform_from and user_platform_to cannot be the same")
		}
	}

	return nil
}

func (s *Service) validatePocket(ctx context.Context, userID primitive.ObjectID, pocketFrom, pocketTo *primitive.ObjectID, amount float64) error {
	// Check pocket from
	if pocketFrom != nil {
		pocket, err := s.pocketRepo.GetPocketByID(ctx, *pocketFrom)
		if err != nil {
			return errors.New("pocket_from not found")
		}

		// Check ownership
		if pocket.UserID != userID {
			return errors.New("unauthorized: pocket_from does not belong to user")
		}

		// Check locked
		if pocket.IsLocked {
			return errors.New("pocket_from is locked")
		}

		// Check active
		if !pocket.IsActive {
			return errors.New("pocket_from is not active")
		}

		// Check sufficient balance
		if utils.Decimal128ToFloat64(pocket.Balance) < amount {
			return errors.New("insufficient pocket balance")
		}
	}

	// Check pocket to
	if pocketTo != nil {
		pocket, err := s.pocketRepo.GetPocketByID(ctx, *pocketTo)
		if err != nil {
			return errors.New("pocket_to not found")
		}

		// Check ownership
		if pocket.UserID != userID {
			return errors.New("unauthorized: pocket_to does not belong to user")
		}

		// Check locked
		if pocket.IsLocked {
			return errors.New("pocket_to is locked")
		}

		// Check active
		if !pocket.IsActive {
			return errors.New("pocket_to is not active")
		}
	}

	return nil
}

func (s *Service) validateUserPlatform(ctx context.Context, userID primitive.ObjectID, userPlatformFrom, userPlatformTo *primitive.ObjectID, amount float64) error {
	if userPlatformFrom != nil {
		userPlatform, err := s.userPlatformRepo.GetUserPlatformByID(ctx, *userPlatformFrom)
		if err != nil {
			return errors.New("user_platform_from not found")
		}
		if userPlatform.UserID != userID {
			return errors.New("unauthorized: user_platform_from does not belong to user")
		}
		if !userPlatform.IsActive {
			return errors.New("user_platform_from is not active")
		}

		if utils.Decimal128ToFloat64(userPlatform.Balance) < amount {
			return errors.New("insufficient user platform balance")
		}
	}

	if userPlatformTo != nil {
		userPlatform, err := s.userPlatformRepo.GetUserPlatformByID(ctx, *userPlatformTo)
		if err != nil {
			return errors.New("user_platform_to not found")
		}
		if userPlatform.UserID != userID {
			return errors.New("unauthorized: user_platform_to does not belong to user")
		}
		if !userPlatform.IsActive {
			return errors.New("user_platform_to is not active")
		}
	}

	return nil
}

func (s *Service) GetTransactionByID(ctx context.Context, userID string, transactionID string) (*Transaction, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	txObjID, err := primitive.ObjectIDFromHex(transactionID)
	if err != nil {
		return nil, errors.New("invalid transaction id")
	}

	transaction, err := s.repo.GetTransactionByID(ctx, txObjID)
	if err != nil {
		return nil, err
	}

	if transaction.UserID != userObjID {
		return nil, errors.New("unauthorized")
	}

	return transaction, nil
}

func (s *Service) GetUserTransactions(ctx context.Context, userID string, limit int64, skip int64) ([]*Transaction, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	return s.repo.GetTransactionsByUserID(ctx, userObjID, limit, skip)
}

func (s *Service) GetUserTransactionsWithSort(ctx context.Context, userID string, txType *string, search *string, page int64, pageSize int64, sortBy string, sortOrder string) ([]*dto.TransactionResponse, int64, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, 0, errors.New("invalid user id")
	}

	return s.repo.GetTransactionsByUserIDWithSort(ctx, userObjID, txType, search, page, pageSize, sortBy, sortOrder)
}

func (s *Service) GetPocketTransactions(ctx context.Context, userID string, pocketID string, limit int64, skip int64) ([]*Transaction, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	pocketObjID, err := primitive.ObjectIDFromHex(pocketID)
	if err != nil {
		return nil, errors.New("invalid pocket id")
	}

	pocket, err := s.pocketRepo.GetPocketByID(ctx, pocketObjID)
	if err != nil {
		return nil, errors.New("pocket not found")
	}

	if pocket.UserID != userObjID {
		return nil, errors.New("unauthorized")
	}

	return s.repo.GetTransactionsByPocketID(ctx, pocketObjID, limit, skip)
}

func (s *Service) GetPocketTransactionsWithSort(ctx context.Context, userID string, pocketID string, page int64, pageSize int64, sortBy string, sortOrder string) ([]*Transaction, int64, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, 0, errors.New("invalid user id")
	}

	pocketObjID, err := primitive.ObjectIDFromHex(pocketID)
	if err != nil {
		return nil, 0, errors.New("invalid pocket id")
	}

	pocket, err := s.pocketRepo.GetPocketByID(ctx, pocketObjID)
	if err != nil {
		return nil, 0, errors.New("pocket not found")
	}

	if pocket.UserID != userObjID {
		return nil, 0, errors.New("unauthorized")
	}

	return s.repo.GetTransactionsByPocketIDWithSort(ctx, pocketObjID, page, pageSize, sortBy, sortOrder)
}

func (s *Service) DeleteTransaction(ctx context.Context, userID string, transactionID string) error {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user id")
	}

	txObjID, err := primitive.ObjectIDFromHex(transactionID)
	if err != nil {
		return errors.New("invalid transaction id")
	}

	transaction, err := s.repo.GetTransactionByID(ctx, txObjID)
	if err != nil {
		return err
	}

	if transaction.UserID != userObjID {
		return errors.New("unauthorized")
	}

	return s.repo.DeleteTransaction(ctx, txObjID)
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
