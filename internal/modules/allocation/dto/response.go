package dto

import "time"

type AllocationResponse struct {
	ID             string     `json:"id"`
	UserID         string     `json:"user_id"`
	PocketID       *string    `json:"pocket_id,omitempty"`
	UserPlatformID *string    `json:"user_platform_id,omitempty"`
	Priority       int        `json:"priority"`
	AllocationType string     `json:"allocation_type"`
	Nominal        float64    `json:"nominal"`
	IsActive       bool       `json:"is_active"`
	ExecuteDay     *int       `json:"execute_day,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}
