package user_category

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserCategory struct {
	ID              primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	UserID          primitive.ObjectID  `bson:"user_id" json:"user_id"`
	TemplateID      *primitive.ObjectID `bson:"template_id,omitempty" json:"template_id,omitempty"`
	Name            string              `bson:"name" json:"name"`
	TransactionType *TransactionType    `bson:"transaction_type,omitempty" json:"transaction_type,omitempty" enums:"income,expense"`
	ParentID        *primitive.ObjectID `bson:"parent_id,omitempty" json:"parent_id,omitempty"`
	Description     *string             `bson:"description,omitempty" json:"description,omitempty"`
	Icon            *string             `bson:"icon,omitempty" json:"icon,omitempty"`
	Color           *string             `bson:"color,omitempty" json:"color,omitempty"`
	IsDefault       bool                `bson:"is_default" json:"is_default"`
	IsDeleted       bool                `bson:"is_deleted" json:"is_deleted"`
	CreatedAt       time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time           `bson:"updated_at" json:"updated_at"`
	DeletedAt       *time.Time          `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

type TransactionType string

const (
	TransactionIncome  TransactionType = "income"
	TransactionExpense TransactionType = "expense"
)
