package dto

type CreateCategoryRequest struct {
	Name            string  `json:"name" validate:"required,min=1,max=255"`
	Type            string  `json:"type" validate:"required,oneof=transaction pocket"`
	TransactionType *string `json:"transaction_type,omitempty" validate:"omitempty,oneof=income expense"`
	IsDefault       bool    `json:"is_default"`
	ParentID        string  `json:"parent_id" validate:"omitempty,len=24,hexadecimal"`
	Description     *string `json:"description,omitempty" validate:"omitempty,max=500"`
	Icon            *string `json:"icon,omitempty" validate:"omitempty,max=100"`
	Color           *string `json:"color,omitempty" validate:"omitempty,max=50"`
}

type UpdateCategoryRequest struct {
	Name            string  `json:"name" validate:"omitempty,min=1,max=255"`
	Type            string  `json:"type" validate:"omitempty,oneof=transaction pocket"`
	TransactionType *string `json:"transaction_type,omitempty" validate:"omitempty,oneof=income expense"`
	IsDefault       bool    `json:"is_default"`
	ParentID        string  `json:"parent_id" validate:"omitempty,len=24,hexadecimal"`
	Description     *string `json:"description,omitempty" validate:"omitempty,max=500"`
	Icon            *string `json:"icon,omitempty" validate:"omitempty,max=100"`
	Color           *string `json:"color,omitempty" validate:"omitempty,max=50"`
}
