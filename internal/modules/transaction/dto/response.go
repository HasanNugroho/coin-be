package dto

import "time"

type TransactionResponse struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	Type       string     `json:"type"`
	Amount     float64    `json:"amount"`
	PocketFrom *string    `json:"pocket_from,omitempty"`
	PocketTo   *string    `json:"pocket_to,omitempty"`
	CategoryID *string    `json:"category_id,omitempty"`
	PlatformID *string    `json:"platform_id,omitempty"`
	Note       *string    `json:"note,omitempty"`
	Date       time.Time  `json:"date"`
	Ref        *string    `json:"ref,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}
