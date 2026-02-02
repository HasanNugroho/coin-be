package dto

type CreateTransactionRequest struct {
	Type               string  `json:"type" validate:"required,oneof=income expense transfer"`
	Amount             float64 `json:"amount" validate:"required,gt=0"`
	PocketFromID       string  `json:"pocket_from_id" validate:"omitempty,len=24,hexadecimal"`
	PocketToID         string  `json:"pocket_to_id" validate:"omitempty,len=24,hexadecimal"`
	UserPlatformFromID string  `json:"user_platform_from_id" validate:"omitempty,len=24,hexadecimal"`
	UserPlatformToID   string  `json:"user_platform_to_id" validate:"omitempty,len=24,hexadecimal"`
	CategoryID         string  `json:"category_id" validate:"omitempty,len=24,hexadecimal"`
	Note               string  `json:"note" validate:"omitempty,max=500"`
	Date               string  `json:"date" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	Ref                string  `json:"ref" validate:"omitempty,max=100"`
}
