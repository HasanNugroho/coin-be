package dto

type CreatePlatformRequest struct {
	Name     string `json:"name" validate:"required,min=1,max=255"`
	Type     string `json:"type" validate:"required,oneof=BANK E_WALLET CASH ATM"`
	IsActive bool   `json:"is_active"`
}

type UpdatePlatformRequest struct {
	Name     string `json:"name" validate:"omitempty,min=1,max=255"`
	Type     string `json:"type" validate:"omitempty,oneof=BANK E_WALLET CASH ATM"`
	IsActive *bool  `json:"is_active"`
}
