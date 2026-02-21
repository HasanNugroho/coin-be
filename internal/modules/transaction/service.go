package transaction

import (
	"context"
	"errors"
	"time"

	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/internal/modules/dashboard"
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
	dashboardService *dashboard.Service
}

func NewService(r *Repository, pr *pocket.Repository, upr *user_platform.UserPlatformRepository, ds *dashboard.Service) *Service {
	return &Service{
		repo:             r,
		pocketRepo:       pr,
		userPlatformRepo: upr,
		balanceProcessor: NewBalanceProcessor(pr, upr),
		dashboardService: ds,
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

	// Trigger daily summary recalculation if the transaction date is in the past (date has passed)
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if date.Before(today) {
		go s.dashboardService.GenerateDailySummary(context.Background(), userObjID, date)
	}

	return transaction, nil
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

	// Process balance updates (revert)
	if err := s.balanceProcessor.RevertTransaction(ctx, transaction.Type, transaction.Amount, transaction.PocketFromID, transaction.PocketToID, transaction.UserPlatformFromID, transaction.UserPlatformToID); err != nil {
		return err
	}

	err = s.repo.DeleteTransaction(ctx, txObjID)
	if err != nil {
		return err
	}

	// Trigger daily summary recalculation if the transaction date is in the past (date has passed)
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if transaction.Date.Before(today) {
		go s.dashboardService.GenerateDailySummary(context.Background(), transaction.UserID, transaction.Date)
	}

	return nil
}

func (s *Service) UpdateTransaction(ctx context.Context, userID string, transactionID string, req *dto.UpdateTransactionRequest) (*Transaction, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	txObjID, err := primitive.ObjectIDFromHex(transactionID)
	if err != nil {
		return nil, errors.New("invalid transaction id")
	}

	// 1. Fetch existing transaction
	oldTx, err := s.repo.GetTransactionByID(ctx, txObjID)
	if err != nil {
		return nil, err
	}

	if oldTx.UserID != userObjID {
		return nil, errors.New("unauthorized")
	}

	// 2. Validate request
	if !IsValidTransactionType(req.Type) {
		return nil, errors.New("invalid transaction type")
	}

	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	newDate, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		return nil, errors.New("invalid date format")
	}

	var newPocketFrom *primitive.ObjectID
	var newPocketTo *primitive.ObjectID
	var newUserPlatformFrom *primitive.ObjectID
	var newUserPlatformTo *primitive.ObjectID

	if req.PocketFromID != "" {
		id, _ := primitive.ObjectIDFromHex(req.PocketFromID)
		newPocketFrom = &id
	}
	if req.PocketToID != "" {
		id, _ := primitive.ObjectIDFromHex(req.PocketToID)
		newPocketTo = &id
	}
	if req.UserPlatformFromID != "" {
		id, _ := primitive.ObjectIDFromHex(req.UserPlatformFromID)
		newUserPlatformFrom = &id
	}
	if req.UserPlatformToID != "" {
		id, _ := primitive.ObjectIDFromHex(req.UserPlatformToID)
		newUserPlatformTo = &id
	}

	var newCategoryID *primitive.ObjectID
	if req.CategoryID != "" {
		id, _ := primitive.ObjectIDFromHex(req.CategoryID)
		newCategoryID = &id
	}

	// Validate generic rules
	if err := s.validateTransactionRules(ctx, req.Type, userObjID, newPocketFrom, newPocketTo, newUserPlatformFrom, newUserPlatformTo); err != nil {
		return nil, err
	}

	// 3. Revert old balances
	if err := s.balanceProcessor.RevertTransaction(ctx, oldTx.Type, oldTx.Amount, oldTx.PocketFromID, oldTx.PocketToID, oldTx.UserPlatformFromID, oldTx.UserPlatformToID); err != nil {
		return nil, err
	}

	// 4. Validate new ownership and balance sufficiency (after reversion)
	if err := s.validatePocket(ctx, userObjID, newPocketFrom, newPocketTo, req.Amount); err != nil {
		// Rollback reversion if validation fails
		s.balanceProcessor.ProcessTransaction(ctx, oldTx.Type, oldTx.Amount, oldTx.PocketFromID, oldTx.PocketToID, oldTx.UserPlatformFromID, oldTx.UserPlatformToID)
		return nil, err
	}

	if err := s.validateUserPlatform(ctx, userObjID, newUserPlatformFrom, newUserPlatformTo, req.Amount); err != nil {
		// Rollback reversion if validation fails
		s.balanceProcessor.ProcessTransaction(ctx, oldTx.Type, oldTx.Amount, oldTx.PocketFromID, oldTx.PocketToID, oldTx.UserPlatformFromID, oldTx.UserPlatformToID)
		return nil, err
	}

	// 5. Apply new balances
	if err := s.balanceProcessor.ProcessTransaction(ctx, req.Type, req.Amount, newPocketFrom, newPocketTo, newUserPlatformFrom, newUserPlatformTo); err != nil {
		// Rollback reversion if apply fails
		s.balanceProcessor.ProcessTransaction(ctx, oldTx.Type, oldTx.Amount, oldTx.PocketFromID, oldTx.PocketToID, oldTx.UserPlatformFromID, oldTx.UserPlatformToID)
		return nil, err
	}

	// 6. Update transaction record
	updatedTx := &Transaction{
		ID:                 txObjID,
		UserID:             userObjID,
		Type:               req.Type,
		Amount:             req.Amount,
		PocketFromID:       newPocketFrom,
		PocketToID:         newPocketTo,
		UserPlatformFromID: newUserPlatformFrom,
		UserPlatformToID:   newUserPlatformTo,
		CategoryID:         newCategoryID,
		Note:               stringPtr(req.Note),
		Date:               newDate,
		Ref:                stringPtr(req.Ref),
		CreatedAt:          oldTx.CreatedAt,
	}

	if err := s.repo.UpdateTransaction(ctx, txObjID, updatedTx); err != nil {
		// FATAL: Balance already updated but DB update failed.
		// In a real system we'd use a transaction if MongoDB supported it across collections easily here.
		return nil, err
	}

	// 7. Trigger daily summary recalculation for both old and new dates if they are in the past
	loc := time.Local // atau time.LoadLocation("Asia/Jakarta")

	now := time.Now().In(loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	oldDate := oldTx.Date.In(loc)
	newDateLoc := newDate.In(loc)

	// trigger untuk tanggal lama kalau itu di masa lalu
	if oldDate.Before(today) {
		go s.dashboardService.GenerateDailySummary(context.Background(), userObjID, oldDate)
	}

	// kalau tanggal berubah (beda hari)
	if oldDate.Year() != newDateLoc.Year() ||
		oldDate.Month() != newDateLoc.Month() ||
		oldDate.Day() != newDateLoc.Day() {

		// trigger untuk tanggal baru kalau di masa lalu
		if newDateLoc.Before(today) {
			go s.dashboardService.GenerateDailySummary(context.Background(), userObjID, newDateLoc)
		}
	}

	return updatedTx, nil
}

func (s *Service) validateTransactionRules(
	ctx context.Context,
	txType string,
	userID primitive.ObjectID,
	pocketFrom *primitive.ObjectID,
	pocketTo *primitive.ObjectID,
	userPlatformFrom *primitive.ObjectID,
	userPlatformTo *primitive.ObjectID,
) error {
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

func (s *Service) validatePocket(
	ctx context.Context,
	userID primitive.ObjectID,
	pocketFrom *primitive.ObjectID,
	pocketTo *primitive.ObjectID,
	amount float64,
) error {
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

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
