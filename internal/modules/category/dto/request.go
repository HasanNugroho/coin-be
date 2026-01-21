package dto

type CreateCategoryRequest struct {
	Name  string `json:"name" binding:"required"`
	Type  string `json:"type" binding:"required,oneof=income expense"`
	Icon  string `json:"icon"`
	Color string `json:"color"`
}

type UpdateCategoryRequest struct {
	Name  string `json:"name"`
	Icon  string `json:"icon"`
	Color string `json:"color"`
}
