package dto

type CreateUserPlatformRequest struct {
	PlatformID string  `json:"platform_id" validate:"required,len=24,hexadecimal"`
	AliasName  *string `json:"alias_name" validate:"omitempty,min=1,max=255"`
}

type UpdateUserPlatformRequest struct {
	AliasName *string `json:"alias_name" validate:"omitempty,min=1,max=255"`
	IsActive  *bool   `json:"is_active"`
}
