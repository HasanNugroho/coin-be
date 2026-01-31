package dto

type CreatePocketTemplateRequest struct {
	Name            string `json:"name" validate:"required,min=1,max=255"`
	Type            string `json:"type" validate:"required,oneof=main saving allocation"`
	CategoryID      string `json:"category_id" validate:"required,len=24,hexadecimal"`
	IsDefault       bool   `json:"is_default"`
	IsActive        bool   `json:"is_active"`
	Order           int    `json:"order" validate:"min=0,max=10000"`
	Icon            string `json:"icon" validate:"omitempty,max=100"`
	IconColor       string `json:"icon_color" validate:"omitempty,max=50"`
	BackgroundColor string `json:"background_color" validate:"omitempty,max=50"`
}

type UpdatePocketTemplateRequest struct {
	Name            string `json:"name" validate:"omitempty,min=1,max=255"`
	Type            string `json:"type" validate:"omitempty,oneof=main saving allocation"`
	CategoryID      string `json:"category_id" validate:"omitempty,len=24,hexadecimal"`
	IsDefault       bool   `json:"is_default"`
	IsActive        bool   `json:"is_active"`
	Order           int    `json:"order" validate:"omitempty,min=0,max=10000"`
	Icon            string `json:"icon" validate:"omitempty,max=100"`
	IconColor       string `json:"icon_color" validate:"omitempty,max=50"`
	BackgroundColor string `json:"background_color" validate:"omitempty,max=50"`
}
