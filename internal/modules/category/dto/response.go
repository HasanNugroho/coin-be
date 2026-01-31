package dto

import (
	"time"
)

type CategoryResponse struct {
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	Type            string     `json:"type"`
	TransactionType *string    `json:"transaction_type,omitempty"`
	IsDefault       bool       `json:"is_default"`
	Color           *string    `json:"color,omitempty"`
	Icon            *string    `json:"icon,omitempty"`
	Description     *string    `json:"description,omitempty"`
	ParentID        *string    `json:"parent_id,omitempty"`
	UserID          *string    `json:"user_id,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}
