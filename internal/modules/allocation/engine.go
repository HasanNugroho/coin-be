package allocation

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AllocationEngine struct {
	allocationRepo *Repository
}

func NewAllocationEngine(allocationRepo *Repository) *AllocationEngine {
	return &AllocationEngine{
		allocationRepo: allocationRepo,
	}
}

type DistributionResult struct {
	AllocationID    primitive.ObjectID `json:"allocation_id"`
	AllocationName  string             `json:"allocation_name"`
	AllocatedAmount float64            `json:"allocated_amount"`
	Percentage      float64            `json:"percentage"`
	Priority        int                `json:"priority"`
}

type DistributionSummary struct {
	TotalIncome      float64              `json:"total_income"`
	TotalDistributed float64              `json:"total_distributed"`
	FreeCash         float64              `json:"free_cash"`
	Distributions    []DistributionResult `json:"distributions"`
}

func (e *AllocationEngine) DistributeIncome(ctx context.Context, userID primitive.ObjectID, transactionID primitive.ObjectID, income float64) (*DistributionSummary, error) {
	allocations, err := e.allocationRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active allocations: %w", err)
	}

	remaining := income
	var distributions []DistributionResult

	for _, allocation := range allocations {
		if allocation.TargetAmount != nil && allocation.CurrentAmount >= *allocation.TargetAmount {
			continue
		}

		ideal := income * allocation.Percentage / 100

		allocated := ideal
		if allocated > remaining {
			allocated = remaining
		}

		if allocation.TargetAmount != nil {
			maxAllowed := *allocation.TargetAmount - allocation.CurrentAmount
			if allocated > maxAllowed {
				allocated = maxAllowed
			}
		}

		if allocated > 0 {
			log := &AllocationLog{
				UserID:          userID,
				AllocationID:    allocation.ID,
				TransactionID:   transactionID,
				IncomeAmount:    income,
				AllocatedAmount: allocated,
				Percentage:      allocation.Percentage,
				Priority:        allocation.Priority,
			}

			if err := e.allocationRepo.CreateLog(ctx, log); err != nil {
				return nil, fmt.Errorf("failed to create allocation log: %w", err)
			}

			if err := e.allocationRepo.UpdateCurrentAmount(ctx, allocation.ID, allocated); err != nil {
				return nil, fmt.Errorf("failed to update allocation amount: %w", err)
			}

			distributions = append(distributions, DistributionResult{
				AllocationID:    allocation.ID,
				AllocationName:  allocation.Name,
				AllocatedAmount: allocated,
				Percentage:      allocation.Percentage,
				Priority:        allocation.Priority,
			})

			remaining -= allocated
		}

		if remaining <= 0 {
			break
		}
	}

	totalDistributed := income - remaining

	return &DistributionSummary{
		TotalIncome:      income,
		TotalDistributed: totalDistributed,
		FreeCash:         remaining,
		Distributions:    distributions,
	}, nil
}

func (e *AllocationEngine) ValidateExpense(ctx context.Context, userID primitive.ObjectID, amount float64, allocationID *primitive.ObjectID) error {
	if allocationID != nil {
		allocation, err := e.allocationRepo.GetByID(ctx, *allocationID)
		if err != nil {
			return fmt.Errorf("allocation not found: %w", err)
		}

		if allocation.CurrentAmount < amount {
			return fmt.Errorf("insufficient balance in allocation %s: available %.2f, required %.2f",
				allocation.Name, allocation.CurrentAmount, amount)
		}
	}

	return nil
}

func (e *AllocationEngine) ProcessExpense(ctx context.Context, userID primitive.ObjectID, amount float64, allocationID *primitive.ObjectID) error {
	if allocationID != nil {
		if err := e.allocationRepo.UpdateCurrentAmount(ctx, *allocationID, -amount); err != nil {
			return fmt.Errorf("failed to update allocation amount: %w", err)
		}
	}

	return nil
}
