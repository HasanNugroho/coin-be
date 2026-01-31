package dto

type CreateCategoryRequest struct {
	Name            string  `json:"name" binding:"required"`
	Type            string  `json:"type" binding:"required,oneof=transaction pocket"`
	TransactionType *string `json:"transaction_type,omitempty" binding:"omitempty,oneof=income expense"`
	IsDefault       bool    `json:"is_default"`
	ParentID        string  `json:"parent_id"`
	Description     *string `json:"description,omitempty"`
	Icon            *string `json:"icon,omitempty"`
	Color           *string `json:"color,omitempty"`
}

type UpdateCategoryRequest struct {
	Name            string  `json:"name"`
	Type            string  `json:"type" binding:"omitempty,oneof=transaction pocket"`
	TransactionType *string `json:"transaction_type,omitempty" binding:"omitempty,oneof=income expense"`
	IsDefault       bool    `json:"is_default"`
	ParentID        string  `json:"parent_id"`
	Description     *string `json:"description,omitempty"`
	Icon            *string `json:"icon,omitempty"`
	Color           *string `json:"color,omitempty"`
}
