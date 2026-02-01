package dto

type CreateTransactionRequest struct {
	Type       string  `json:"type" validate:"required,oneof=income expense transfer dp withdraw"`
	Amount     float64 `json:"amount" validate:"required,gt=0"`
	PocketFrom string  `json:"pocket_from" validate:"omitempty,len=24,hexadecimal"`
	PocketTo   string  `json:"pocket_to" validate:"omitempty,len=24,hexadecimal"`
	CategoryID string  `json:"category_id" validate:"omitempty,len=24,hexadecimal"`
	PlatformID string  `json:"platform_id" validate:"omitempty,len=24,hexadecimal"`
	Note       string  `json:"note" validate:"omitempty,max=500"`
	Date       string  `json:"date" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	Ref        string  `json:"ref" validate:"omitempty,max=100"`
}
