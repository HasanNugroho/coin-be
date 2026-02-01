package transaction

import (
	"context"
	"errors"
	"time"

	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/internal/modules/pocket"
	"github.com/HasanNugroho/coin-be/internal/modules/transaction/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo       *Repository
	pocketRepo *pocket.Repository
}

func NewService(r *Repository, pr *pocket.Repository) *Service {
	return &Service{
		repo:       r,
		pocketRepo: pr,
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

	if req.PocketFrom != "" {
		pocketFromID, err := primitive.ObjectIDFromHex(req.PocketFrom)
		if err != nil {
			return nil, errors.New("invalid pocket_from id")
		}
		pocketFrom = &pocketFromID
	}

	if req.PocketTo != "" {
		pocketToID, err := primitive.ObjectIDFromHex(req.PocketTo)
		if err != nil {
			return nil, errors.New("invalid pocket_to id")
		}
		pocketTo = &pocketToID
	}

	var categoryID *primitive.ObjectID
	if req.CategoryID != "" {
		catID, err := primitive.ObjectIDFromHex(req.CategoryID)
		if err != nil {
			return nil, errors.New("invalid category id")
		}
		categoryID = &catID
	}

	var platformID *primitive.ObjectID
	if req.PlatformID != "" {
		platID, err := primitive.ObjectIDFromHex(req.PlatformID)
		if err != nil {
			return nil, errors.New("invalid platform id")
		}
		platformID = &platID
	}

	if err := s.validateTransactionRules(ctx, req.Type, userObjID, pocketFrom, pocketTo); err != nil {
		return nil, err
	}

	if err := s.validatePocketOwnership(ctx, userObjID, pocketFrom, pocketTo); err != nil {
		return nil, err
	}

	if err := s.validatePocketStatus(ctx, pocketFrom, pocketTo); err != nil {
		return nil, err
	}

	if err := s.validateSufficientBalance(ctx, req.Type, pocketFrom, req.Amount); err != nil {
		return nil, err
	}

	transaction := &Transaction{
		UserID:     userObjID,
		Type:       req.Type,
		Amount:     req.Amount,
		PocketFrom: pocketFrom,
		PocketTo:   pocketTo,
		CategoryID: categoryID,
		PlatformID: platformID,
		Note:       stringPtr(req.Note),
		Date:       date,
		Ref:        stringPtr(req.Ref),
	}

	if err := s.repo.CreateTransaction(ctx, transaction); err != nil {
		return nil, err
	}

	if err := s.updatePocketBalances(ctx, req.Type, pocketFrom, pocketTo, req.Amount); err != nil {
		s.repo.DeleteTransaction(ctx, transaction.ID)
		return nil, errors.New("failed to update pocket balances")
	}

	return transaction, nil
}

func (s *Service) validateTransactionRules(ctx context.Context, txType string, userID primitive.ObjectID, pocketFrom, pocketTo *primitive.ObjectID) error {
	switch txType {
	case string(TypeIncome):
		if pocketTo == nil {
			return errors.New("pocket_to is required for INCOME transactions")
		}
		if pocketFrom != nil {
			return errors.New("pocket_from must be null for INCOME transactions")
		}

	case string(TypeExpense):
		if pocketFrom == nil {
			return errors.New("pocket_from is required for EXPENSE transactions")
		}
		if pocketTo != nil {
			return errors.New("pocket_to must be null for EXPENSE transactions")
		}

	case string(TypeTransfer):
		if pocketFrom == nil || pocketTo == nil {
			return errors.New("both pocket_from and pocket_to are required for TRANSFER transactions")
		}
		if pocketFrom.Hex() == pocketTo.Hex() {
			return errors.New("pocket_from and pocket_to cannot be the same")
		}

	case string(TypeDebtPayment):
		if pocketFrom == nil {
			return errors.New("pocket_from is required for DEBT_PAYMENT transactions")
		}

	case string(TypeWithdraw):
		if pocketFrom == nil {
			return errors.New("pocket_from is required for WITHDRAW transactions")
		}
		if pocketTo != nil {
			return errors.New("pocket_to must be null for WITHDRAW transactions")
		}
	}

	return nil
}

func (s *Service) validatePocketOwnership(ctx context.Context, userID primitive.ObjectID, pocketFrom, pocketTo *primitive.ObjectID) error {
	if pocketFrom != nil {
		pocket, err := s.pocketRepo.GetPocketByID(ctx, *pocketFrom)
		if err != nil {
			return errors.New("pocket_from not found")
		}
		if pocket.UserID != userID {
			return errors.New("unauthorized: pocket_from does not belong to user")
		}
	}

	if pocketTo != nil {
		pocket, err := s.pocketRepo.GetPocketByID(ctx, *pocketTo)
		if err != nil {
			return errors.New("pocket_to not found")
		}
		if pocket.UserID != userID {
			return errors.New("unauthorized: pocket_to does not belong to user")
		}
	}

	return nil
}

func (s *Service) validatePocketStatus(ctx context.Context, pocketFrom, pocketTo *primitive.ObjectID) error {
	if pocketFrom != nil {
		pocket, err := s.pocketRepo.GetPocketByID(ctx, *pocketFrom)
		if err != nil {
			return err
		}
		if pocket.IsLocked {
			return errors.New("pocket_from is locked")
		}
		if !pocket.IsActive {
			return errors.New("pocket_from is not active")
		}
	}

	if pocketTo != nil {
		pocket, err := s.pocketRepo.GetPocketByID(ctx, *pocketTo)
		if err != nil {
			return err
		}
		if pocket.IsLocked {
			return errors.New("pocket_to is locked")
		}
		if !pocket.IsActive {
			return errors.New("pocket_to is not active")
		}
	}

	return nil
}

func (s *Service) validateSufficientBalance(ctx context.Context, txType string, pocketFrom *primitive.ObjectID, amount float64) error {
	if pocketFrom == nil {
		return nil
	}

	pocket, err := s.pocketRepo.GetPocketByID(ctx, *pocketFrom)
	if err != nil {
		return err
	}

	if utils.Decimal128ToFloat64(pocket.Balance) < amount {
		return errors.New("insufficient balance")
	}

	return nil
}

func (s *Service) updatePocketBalances(ctx context.Context, txType string, pocketFrom, pocketTo *primitive.ObjectID, amount float64) error {
	switch txType {
	case string(TypeIncome):
		if pocketTo != nil {
			pocketData, err := s.pocketRepo.GetPocketByID(ctx, *pocketTo)
			if err != nil {
				return err
			}

			pocketData.Balance = utils.AddDecimal128(pocketData.Balance, amount)

			if err := s.pocketRepo.UpdatePocket(ctx, *pocketTo, pocketData); err != nil {
				return err
			}
		}

	case string(TypeExpense):
		if pocketFrom != nil {
			pocketData, err := s.pocketRepo.GetPocketByID(ctx, *pocketFrom)
			if err != nil {
				return err
			}
			pocketData.Balance = utils.AddDecimal128(pocketData.Balance, -amount)
			if err := s.pocketRepo.UpdatePocket(ctx, *pocketFrom, pocketData); err != nil {
				return err
			}
		}

	case string(TypeTransfer):
		if pocketFrom != nil {
			pocketFromData, err := s.pocketRepo.GetPocketByID(ctx, *pocketFrom)
			if err != nil {
				return err
			}
			pocketFromData.Balance = utils.AddDecimal128(pocketFromData.Balance, -amount)
			if err := s.pocketRepo.UpdatePocket(ctx, *pocketFrom, pocketFromData); err != nil {
				return err
			}
		}

		if pocketTo != nil {
			pocketToData, err := s.pocketRepo.GetPocketByID(ctx, *pocketTo)
			if err != nil {
				return err
			}
			pocketToData.Balance = utils.AddDecimal128(pocketToData.Balance, amount)
			if err := s.pocketRepo.UpdatePocket(ctx, *pocketTo, pocketToData); err != nil {
				return err
			}
		}

	case string(TypeDebtPayment):
		if pocketFrom != nil {
			pocketData, err := s.pocketRepo.GetPocketByID(ctx, *pocketFrom)
			if err != nil {
				return err
			}
			pocketData.Balance = utils.AddDecimal128(pocketData.Balance, -amount)
			if err := s.pocketRepo.UpdatePocket(ctx, *pocketFrom, pocketData); err != nil {
				return err
			}
		}

	case string(TypeWithdraw):
		if pocketFrom != nil {
			pocketData, err := s.pocketRepo.GetPocketByID(ctx, *pocketFrom)
			if err != nil {
				return err
			}
			pocketData.Balance = utils.AddDecimal128(pocketData.Balance, -amount)
			if err := s.pocketRepo.UpdatePocket(ctx, *pocketFrom, pocketData); err != nil {
				return err
			}
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
