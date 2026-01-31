package pocket_template

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PocketTemplate struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	Name       string              `bson:"name" json:"name"`
	Type       string              `bson:"type" json:"type" enums:"main,saving,allocation"`
	CategoryID *primitive.ObjectID `bson:"category_id,omitempty" json:"category_id,omitempty"`

	IsDefault bool `bson:"is_default" json:"is_default"`
	IsActive  bool `bson:"is_active" json:"is_active"`

	Order int `bson:"order" json:"order"`

	Icon            *string `bson:"icon,omitempty" json:"icon,omitempty"`
	IconColor       *string `bson:"icon_color,omitempty" json:"icon_color,omitempty"`
	BackgroundColor *string `bson:"background_color,omitempty" json:"background_color,omitempty"`

	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

type PocketTemplateType string

const (
	TypeMain       PocketTemplateType = "main"
	TypeSaving     PocketTemplateType = "saving"
	TypeAllocation PocketTemplateType = "allocation"
)
