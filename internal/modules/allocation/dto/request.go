package dto

type CreateAllocationRequest struct {
	PocketID       string  `json:"pocket_id" validate:"omitempty,len=24,hexadecimal"`
	UserPlatformID string  `json:"user_platform_id" validate:"omitempty,len=24,hexadecimal"`
	Priority       int     `json:"priority" validate:"required,min=1,max=3"`
	AllocationType string  `json:"allocation_type" validate:"required,oneof=PERCENTAGE NOMINAL"`
	Nominal        float64 `json:"nominal" validate:"required,gt=0"`
	ExecuteDay     *int    `json:"execute_day" validate:"omitempty,min=1,max=31"`
}

type UpdateAllocationRequest struct {
	PocketID       string   `json:"pocket_id" validate:"omitempty,len=24,hexadecimal"`
	UserPlatformID string   `json:"user_platform_id" validate:"omitempty,len=24,hexadecimal"`
	Priority       *int     `json:"priority" validate:"omitempty,min=1,max=3"`
	AllocationType string   `json:"allocation_type" validate:"omitempty,oneof=PERCENTAGE NOMINAL"`
	Nominal        *float64 `json:"nominal" validate:"omitempty,gt=0"`
	IsActive       *bool    `json:"is_active"`
	ExecuteDay     *int     `json:"execute_day" validate:"omitempty,min=1,max=31"`
}
