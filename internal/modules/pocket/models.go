package pocket

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Pocket struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	UserID     primitive.ObjectID  `bson:"user_id" json:"user_id"`
	Name       string              `bson:"name" json:"name"`
	Type       string              `bson:"type" json:"type" enums:"main,allocation,saving,debt,system"`
	CategoryID *primitive.ObjectID `bson:"category_id,omitempty" json:"category_id,omitempty"`

	Balance       primitive.Decimal128  `bson:"balance" json:"balance"`
	TargetBalance *primitive.Decimal128 `bson:"target_balance,omitempty" json:"target_balance,omitempty"`

	IsDefault bool `bson:"is_default" json:"is_default"`
	IsActive  bool `bson:"is_active" json:"is_active"`
	IsLocked  bool `bson:"is_locked" json:"is_locked"`

	Icon            string `bson:"icon,omitempty" json:"icon,omitempty"`
	IconColor       string `bson:"icon_color,omitempty" json:"icon_color,omitempty"`
	BackgroundColor string `bson:"background_color,omitempty" json:"background_color,omitempty"`

	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

type PocketType string

const (
	TypeMain       PocketType = "main"
	TypeAllocation PocketType = "allocation"
	TypeSaving     PocketType = "saving"
	TypeDebt       PocketType = "debt"
	TypeSystem     PocketType = "system"
)
