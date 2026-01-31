package dto

type CreatePocketRequest struct {
	Name            string  `json:"name" validate:"required,min=2,max=255"`
	Type            string  `json:"type" validate:"required,oneof=main allocation saving debt"`
	CategoryID      string  `json:"category_id" validate:"omitempty,len=24,hexadecimal"`
	Icon            string  `json:"icon" validate:"omitempty,max=100"`
	IconColor       string  `json:"icon_color" validate:"omitempty,max=50"`
	BackgroundColor string  `json:"background_color" validate:"omitempty,max=50"`
}

type UpdatePocketRequest struct {
	Name            string  `json:"name" validate:"omitempty,min=2,max=255"`
	Type            string  `json:"type" validate:"omitempty,oneof=main allocation saving debt"`
	CategoryID      string  `json:"category_id" validate:"omitempty,len=24,hexadecimal"`
	Icon            string  `json:"icon" validate:"omitempty,max=100"`
	IconColor       string  `json:"icon_color" validate:"omitempty,max=50"`
	BackgroundColor string  `json:"background_color" validate:"omitempty,max=50"`
	IsActive        *bool   `json:"is_active"`
}

type CreateSystemPocketRequest struct {
	Name            string `json:"name" validate:"required,min=2,max=255"`
	CategoryID      string `json:"category_id" validate:"omitempty,len=24,hexadecimal"`
	Icon            string `json:"icon" validate:"omitempty,max=100"`
	IconColor       string `json:"icon_color" validate:"omitempty,max=50"`
	BackgroundColor string `json:"background_color" validate:"omitempty,max=50"`
}
