package dto

import (
	"time"
)

type PocketTemplateResponse struct {
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	Type            string     `json:"type"`
	IsDefault       bool       `json:"is_default"`
	IsActive        bool       `json:"is_active"`
	Order           int        `json:"order"`
	Icon            string     `json:"icon,omitempty"`
	IconColor       string     `json:"icon_color,omitempty"`
	BackgroundColor string     `json:"background_color,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}
