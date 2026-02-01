package dto

type CreateUserCategoryRequest struct {
	Name            string  `json:"name" validate:"required,min=1,max=255"`
	TemplateID      string  `json:"template_id" validate:"omitempty,len=24,hexadecimal"`
	TransactionType *string `json:"transaction_type,omitempty" validate:"omitempty,oneof=income expense"`
	IsDefault       bool    `json:"is_default"`
	ParentID        string  `json:"parent_id" validate:"omitempty,len=24,hexadecimal"`
	Description     *string `json:"description,omitempty" validate:"omitempty,max=500"`
	Icon            *string `json:"icon,omitempty" validate:"omitempty,max=100"`
	Color           *string `json:"color,omitempty" validate:"omitempty,max=50"`
}

type UpdateUserCategoryRequest struct {
	Name            string  `json:"name" validate:"omitempty,min=1,max=255"`
	TemplateID      string  `json:"template_id" validate:"omitempty,len=24,hexadecimal"`
	TransactionType *string `json:"transaction_type,omitempty" validate:"omitempty,oneof=income expense"`
	IsDefault       bool    `json:"is_default"`
	ParentID        string  `json:"parent_id" validate:"omitempty,len=24,hexadecimal"`
	Description     *string `json:"description,omitempty" validate:"omitempty,max=500"`
	Icon            *string `json:"icon,omitempty" validate:"omitempty,max=100"`
	Color           *string `json:"color,omitempty" validate:"omitempty,max=50"`
}
