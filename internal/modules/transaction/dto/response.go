package dto

import "time"

type TransactionResponse struct {
	ID                 string     `bson:"id"                    json:"id"`
	UserID             string     `bson:"user_id"               json:"user_id"`
	Type               string     `bson:"type"                  json:"type"`
	Amount             float64    `bson:"amount"                json:"amount"`
	PocketFromID       *string    `bson:"pocket_from_id"        json:"pocket_from_id,omitempty"`
	PocketFromName     *string    `bson:"pocket_from_name"      json:"pocket_from_name,omitempty"`
	PocketToID         *string    `bson:"pocket_to_id"          json:"pocket_to_id,omitempty"`
	PocketToName       *string    `bson:"pocket_to_name"        json:"pocket_to_name,omitempty"`
	UserPlatformFromID *string    `bson:"user_platform_from_id" json:"user_platform_from_id,omitempty"`
	UserPlatformToID   *string    `bson:"user_platform_to_id"   json:"user_platform_to_id,omitempty"`
	CategoryID         *string    `bson:"category_id"           json:"category_id,omitempty"`
	CategoryName       *string    `bson:"category_name"         json:"category_name,omitempty"`
	Note               *string    `bson:"note"                  json:"note,omitempty"`
	Date               time.Time  `bson:"date"                  json:"date"`
	Ref                *string    `bson:"ref"                   json:"ref,omitempty"`
	CreatedAt          time.Time  `bson:"created_at"            json:"created_at"`
	UpdatedAt          time.Time  `bson:"updated_at"            json:"updated_at"`
	DeletedAt          *time.Time `bson:"deleted_at"            json:"deleted_at,omitempty"`
}
