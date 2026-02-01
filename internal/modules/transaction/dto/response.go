package dto

import "time"

type TransactionResponse struct {
	ID           string     `bson:"id"            json:"id"`
	UserID       string     `bson:"user_id"       json:"user_id"`
	Type         string     `bson:"type"          json:"type"`
	Amount       float64    `bson:"amount"        json:"amount"`
	PocketFrom   *string    `bson:"pocket_from"   json:"pocket_from,omitempty"`
	PocketTo     *string    `bson:"pocket_to"     json:"pocket_to,omitempty"`
	CategoryID   *string    `bson:"category_id"   json:"category_id,omitempty"`
	CategoryName *string    `bson:"category_name" json:"category_name,omitempty"`
	PlatformID   *string    `bson:"platform_id"   json:"platform_id,omitempty"`
	Note         *string    `bson:"note"          json:"note,omitempty"`
	Date         time.Time  `bson:"date"          json:"date"`
	Ref          *string    `bson:"ref"           json:"ref,omitempty"`
	CreatedAt    time.Time  `bson:"created_at"    json:"created_at"`
	UpdatedAt    time.Time  `bson:"updated_at"    json:"updated_at"`
	DeletedAt    *time.Time `bson:"deleted_at"    json:"deleted_at,omitempty"`
}
