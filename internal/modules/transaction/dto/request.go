package dto

import (
	"time"
)

type CreateTransactionRequest struct {
	Type            string    `json:"type" binding:"required,oneof=income expense"`
	Amount          float64   `json:"amount" binding:"required,gt=0"`
	CategoryID      string    `json:"category_id" binding:"required"`
	AllocationID    *string   `json:"allocation_id,omitempty"`
	Description     string    `json:"description"`
	TransactionDate time.Time `json:"transaction_date" binding:"required"`
}

type UpdateTransactionRequest struct {
	Amount          float64   `json:"amount" binding:"gt=0"`
	CategoryID      string    `json:"category_id"`
	AllocationID    *string   `json:"allocation_id,omitempty"`
	Description     string    `json:"description"`
	TransactionDate time.Time `json:"transaction_date"`
}

type FilterTransactionRequest struct {
	Type         string    `form:"type"`
	CategoryID   string    `form:"category_id"`
	AllocationID string    `form:"allocation_id"`
	StartDate    time.Time `form:"start_date"`
	EndDate      time.Time `form:"end_date"`
	Limit        int64     `form:"limit"`
	Skip         int64     `form:"skip"`
}
