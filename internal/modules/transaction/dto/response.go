package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionResponse struct {
	ID              primitive.ObjectID  `json:"id"`
	UserID          primitive.ObjectID  `json:"user_id"`
	Type            string              `json:"type"`
	Amount          float64             `json:"amount"`
	CategoryID      primitive.ObjectID  `json:"category_id"`
	AllocationID    *primitive.ObjectID `json:"allocation_id,omitempty"`
	Description     string              `json:"description"`
	TransactionDate time.Time           `json:"transaction_date"`
	IsDistributed   bool                `json:"is_distributed"`
	CreatedAt       time.Time           `json:"created_at"`
}

type IncomeDistributionResponse struct {
	Transaction   TransactionResponse `json:"transaction"`
	TotalIncome   float64             `json:"total_income"`
	Distributed   float64             `json:"distributed"`
	FreeCash      float64             `json:"free_cash"`
	Distributions []struct {
		AllocationID   primitive.ObjectID `json:"allocation_id"`
		AllocationName string             `json:"allocation_name"`
		Amount         float64            `json:"amount"`
		Percentage     float64            `json:"percentage"`
		Priority       int                `json:"priority"`
	} `json:"distributions"`
}
