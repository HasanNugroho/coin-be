package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AllocationResponse struct {
	ID            primitive.ObjectID `json:"id"`
	UserID        primitive.ObjectID `json:"user_id"`
	Name          string             `json:"name"`
	Priority      int                `json:"priority"`
	Percentage    float64            `json:"percentage"`
	CurrentAmount float64            `json:"current_amount"`
	TargetAmount  *float64           `json:"target_amount,omitempty"`
	IsActive      bool               `json:"is_active"`
	CreatedAt     time.Time          `json:"created_at"`
}

type AllocationLogResponse struct {
	ID              primitive.ObjectID `json:"id"`
	UserID          primitive.ObjectID `json:"user_id"`
	AllocationID    primitive.ObjectID `json:"allocation_id"`
	TransactionID   primitive.ObjectID `json:"transaction_id"`
	IncomeAmount    float64            `json:"income_amount"`
	AllocatedAmount float64            `json:"allocated_amount"`
	Percentage      float64            `json:"percentage"`
	Priority        int                `json:"priority"`
	CreatedAt       time.Time          `json:"created_at"`
}
