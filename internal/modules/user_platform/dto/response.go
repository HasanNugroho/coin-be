package dto

import "time"

type UserPlatformResponse struct {
	ID         string        `json:"id"`
	UserID     string        `json:"user_id"`
	PlatformID string        `json:"platform_id"`
	Platform   *PlatformData `json:"platform,omitempty"`
	AliasName  *string       `json:"alias_name,omitempty"`
	Balance    float64       `json:"balance"`
	IsActive   bool          `json:"is_active"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
	DeletedAt  *time.Time    `json:"deleted_at,omitempty"`
}

type PlatformData struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	IsActive bool   `json:"is_active"`
}

type UserPlatformDropdownResponse struct {
	ID        string        `json:"id"`
	Platform  *PlatformData `json:"platform"`
	AliasName *string       `json:"alias_name,omitempty"`
	Balance   float64       `json:"balance"`
	IsActive  bool          `json:"is_active"`
}
