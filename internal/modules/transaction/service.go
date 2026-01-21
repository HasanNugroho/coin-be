package transaction

import (
	"context"
	"fmt"
	"time"

	"github.com/HasanNugroho/coin-be/internal/modules/allocation"
	"github.com/HasanNugroho/coin-be/internal/modules/transaction/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo   *Repository
	engine *allocation.AllocationEngine
}

func NewService(repo *Repository, engine *allocation.AllocationEngine) *Service {
	return &Service{
		repo:   repo,
		engine: engine,
	}
}

func (s *Service) CreateIncome(ctx context.Context, userID primitive.ObjectID, req *dto.CreateTransactionRequest) (*dto.IncomeDistributionResponse, error) {
	categoryID, err := primitive.ObjectIDFromHex(req.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID: %w", err)
	}

	transaction := &Transaction{
		UserID:          userID,
		Type:            TransactionTypeIncome,
		Amount:          req.Amount,
		CategoryID:      categoryID,
		Description:     req.Description,
		TransactionDate: req.TransactionDate,
		IsDistributed:   true,
	}

	if err := s.repo.Create(ctx, transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	distribution, err := s.engine.DistributeIncome(ctx, userID, transaction.ID, req.Amount)
	if err != nil {
		return nil, fmt.Errorf("failed to distribute income: %w", err)
	}

	response := &dto.IncomeDistributionResponse{
		Transaction: dto.TransactionResponse{
			ID:              transaction.ID,
			UserID:          transaction.UserID,
			Type:            string(transaction.Type),
			Amount:          transaction.Amount,
			CategoryID:      transaction.CategoryID,
			Description:     transaction.Description,
			TransactionDate: transaction.TransactionDate,
			IsDistributed:   transaction.IsDistributed,
			CreatedAt:       transaction.CreatedAt,
		},
		TotalIncome: distribution.TotalIncome,
		Distributed: distribution.TotalDistributed,
		FreeCash:    distribution.FreeCash,
	}

	for _, dist := range distribution.Distributions {
		response.Distributions = append(response.Distributions, struct {
			AllocationID   primitive.ObjectID `json:"allocation_id"`
			AllocationName string             `json:"allocation_name"`
			Amount         float64            `json:"amount"`
			Percentage     float64            `json:"percentage"`
			Priority       int                `json:"priority"`
		}{
			AllocationID:   dist.AllocationID,
			AllocationName: dist.AllocationName,
			Amount:         dist.AllocatedAmount,
			Percentage:     dist.Percentage,
			Priority:       dist.Priority,
		})
	}

	return response, nil
}

func (s *Service) CreateExpense(ctx context.Context, userID primitive.ObjectID, req *dto.CreateTransactionRequest) (*dto.TransactionResponse, error) {
	categoryID, err := primitive.ObjectIDFromHex(req.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID: %w", err)
	}

	var allocationID *primitive.ObjectID
	if req.AllocationID != nil && *req.AllocationID != "" {
		id, err := primitive.ObjectIDFromHex(*req.AllocationID)
		if err != nil {
			return nil, fmt.Errorf("invalid allocation ID: %w", err)
		}
		allocationID = &id
	}

	if err := s.engine.ValidateExpense(ctx, userID, req.Amount, allocationID); err != nil {
		return nil, err
	}

	transaction := &Transaction{
		UserID:          userID,
		Type:            TransactionTypeExpense,
		Amount:          req.Amount,
		CategoryID:      categoryID,
		AllocationID:    allocationID,
		Description:     req.Description,
		TransactionDate: req.TransactionDate,
		IsDistributed:   false,
	}

	if err := s.repo.Create(ctx, transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	if err := s.engine.ProcessExpense(ctx, userID, req.Amount, allocationID); err != nil {
		return nil, fmt.Errorf("failed to process expense: %w", err)
	}

	return &dto.TransactionResponse{
		ID:              transaction.ID,
		UserID:          transaction.UserID,
		Type:            string(transaction.Type),
		Amount:          transaction.Amount,
		CategoryID:      transaction.CategoryID,
		AllocationID:    transaction.AllocationID,
		Description:     transaction.Description,
		TransactionDate: transaction.TransactionDate,
		IsDistributed:   transaction.IsDistributed,
		CreatedAt:       transaction.CreatedAt,
	}, nil
}

func (s *Service) GetTransactions(ctx context.Context, userID primitive.ObjectID, limit, skip int64) ([]*dto.TransactionResponse, error) {
	transactions, err := s.repo.GetByUserID(ctx, userID, limit, skip)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	responses := make([]*dto.TransactionResponse, len(transactions))
	for i, txn := range transactions {
		responses[i] = &dto.TransactionResponse{
			ID:              txn.ID,
			UserID:          txn.UserID,
			Type:            string(txn.Type),
			Amount:          txn.Amount,
			CategoryID:      txn.CategoryID,
			AllocationID:    txn.AllocationID,
			Description:     txn.Description,
			TransactionDate: txn.TransactionDate,
			IsDistributed:   txn.IsDistributed,
			CreatedAt:       txn.CreatedAt,
		}
	}

	return responses, nil
}

func (s *Service) FilterTransactions(ctx context.Context, userID primitive.ObjectID, req *dto.FilterTransactionRequest) ([]*dto.TransactionResponse, error) {
	filter := bson.M{}

	if req.Type != "" {
		filter["type"] = req.Type
	}

	if req.CategoryID != "" {
		categoryID, err := primitive.ObjectIDFromHex(req.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("invalid category ID: %w", err)
		}
		filter["category_id"] = categoryID
	}

	if req.AllocationID != "" {
		allocationID, err := primitive.ObjectIDFromHex(req.AllocationID)
		if err != nil {
			return nil, fmt.Errorf("invalid allocation ID: %w", err)
		}
		filter["allocation_id"] = allocationID
	}

	if !req.StartDate.IsZero() && !req.EndDate.IsZero() {
		filter["transaction_date"] = bson.M{
			"$gte": req.StartDate,
			"$lte": req.EndDate,
		}
	}

	limit := req.Limit
	if limit == 0 {
		limit = 50
	}

	transactions, err := s.repo.Filter(ctx, userID, filter, limit, req.Skip)
	if err != nil {
		return nil, fmt.Errorf("failed to filter transactions: %w", err)
	}

	responses := make([]*dto.TransactionResponse, len(transactions))
	for i, txn := range transactions {
		responses[i] = &dto.TransactionResponse{
			ID:              txn.ID,
			UserID:          txn.UserID,
			Type:            string(txn.Type),
			Amount:          txn.Amount,
			CategoryID:      txn.CategoryID,
			AllocationID:    txn.AllocationID,
			Description:     txn.Description,
			TransactionDate: txn.TransactionDate,
			IsDistributed:   txn.IsDistributed,
			CreatedAt:       txn.CreatedAt,
		}
	}

	return responses, nil
}

func (s *Service) GetTransactionByID(ctx context.Context, id primitive.ObjectID) (*dto.TransactionResponse, error) {
	transaction, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return &dto.TransactionResponse{
		ID:              transaction.ID,
		UserID:          transaction.UserID,
		Type:            string(transaction.Type),
		Amount:          transaction.Amount,
		CategoryID:      transaction.CategoryID,
		AllocationID:    transaction.AllocationID,
		Description:     transaction.Description,
		TransactionDate: transaction.TransactionDate,
		IsDistributed:   transaction.IsDistributed,
		CreatedAt:       transaction.CreatedAt,
	}, nil
}

func (s *Service) UpdateTransaction(ctx context.Context, id primitive.ObjectID, req *dto.UpdateTransactionRequest) (*dto.TransactionResponse, error) {
	transaction, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	if req.Amount > 0 {
		transaction.Amount = req.Amount
	}
	if req.CategoryID != "" {
		categoryID, err := primitive.ObjectIDFromHex(req.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("invalid category ID: %w", err)
		}
		transaction.CategoryID = categoryID
	}
	if req.AllocationID != nil && *req.AllocationID != "" {
		allocationID, err := primitive.ObjectIDFromHex(*req.AllocationID)
		if err != nil {
			return nil, fmt.Errorf("invalid allocation ID: %w", err)
		}
		transaction.AllocationID = &allocationID
	}
	if req.Description != "" {
		transaction.Description = req.Description
	}
	if !req.TransactionDate.IsZero() {
		transaction.TransactionDate = req.TransactionDate
	}

	if err := s.repo.Update(ctx, id, transaction); err != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	return &dto.TransactionResponse{
		ID:              transaction.ID,
		UserID:          transaction.UserID,
		Type:            string(transaction.Type),
		Amount:          transaction.Amount,
		CategoryID:      transaction.CategoryID,
		AllocationID:    transaction.AllocationID,
		Description:     transaction.Description,
		TransactionDate: transaction.TransactionDate,
		IsDistributed:   transaction.IsDistributed,
		CreatedAt:       transaction.CreatedAt,
	}, nil
}

func (s *Service) DeleteTransaction(ctx context.Context, id primitive.ObjectID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}
	return nil
}

func (s *Service) GetTotalIncome(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) (float64, error) {
	return s.repo.GetTotalByUserAndType(ctx, userID, TransactionTypeIncome, startDate, endDate)
}

func (s *Service) GetTotalExpense(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) (float64, error) {
	return s.repo.GetTotalByUserAndType(ctx, userID, TransactionTypeExpense, startDate, endDate)
}
