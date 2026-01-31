package platform

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Platform struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Type      string             `bson:"type" json:"type" enums:"BANK,E_WALLET,CASH,ATM"`
	IsActive  bool               `bson:"is_active" json:"is_active"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time         `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

type PlatformType string

const (
	TypeBank    PlatformType = "BANK"
	TypeEWallet PlatformType = "E_WALLET"
	TypeCash    PlatformType = "CASH"
	TypeATM     PlatformType = "ATM"
)

func IsValidPlatformType(t string) bool {
	switch t {
	case string(TypeBank), string(TypeEWallet), string(TypeCash), string(TypeATM):
		return true
	default:
		return false
	}
}
