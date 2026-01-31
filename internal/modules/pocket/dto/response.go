package dto

import "time"

type PocketResponse struct {
	ID              string     `json:"id"`
	UserID          string     `json:"user_id"`
	Name            string     `json:"name"`
	Type            string     `json:"type"`
	CategoryID      *string    `json:"category_id,omitempty"`
	Balance         float64    `json:"balance"`
	IsDefault       bool       `json:"is_default"`
	IsActive        bool       `json:"is_active"`
	IsLocked        bool       `json:"is_locked"`
	Icon            *string    `json:"icon,omitempty"`
	IconColor       *string    `json:"icon_color,omitempty"`
	BackgroundColor *string    `json:"background_color,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}
