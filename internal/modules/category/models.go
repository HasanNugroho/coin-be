package category

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Category represents a transaction or pocket category
type Category struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// Core
	Name string `bson:"name" json:"name"`

	// Scope
	Type CategoryType `bson:"type" json:"type" enums:"transaction,pocket"`

	// Transaction-only
	TransactionType *TransactionType `bson:"transaction_type,omitempty" json:"transaction_type,omitempty" enums:"income,expense"`

	// Hierarchy
	ParentID *primitive.ObjectID `bson:"parent_id,omitempty" json:"parent_id,omitempty"`

	// Ownership
	UserID *primitive.ObjectID `bson:"user_id,omitempty" json:"user_id,omitempty"` // null = system

	// Metadata
	Description *string `bson:"description,omitempty" json:"description,omitempty"`
	Icon        *string `bson:"icon,omitempty" json:"icon,omitempty"`
	Color       *string `bson:"color,omitempty" json:"color,omitempty"`

	// Flags
	IsDefault bool `bson:"is_default" json:"is_default"`
	IsDeleted bool `bson:"is_deleted" json:"is_deleted"`

	// Audit
	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

// Type constants
const (
	TypeTransaction = "transaction"
	TypePocket      = "pocket"
)

type CategoryType string

const (
	CategoryTypeTransaction CategoryType = "transaction"
	CategoryTypePocket      CategoryType = "pocket"
)

type TransactionType string

const (
	TransactionIncome  TransactionType = "income"
	TransactionExpense TransactionType = "expense"
)

type PocketPurpose string
