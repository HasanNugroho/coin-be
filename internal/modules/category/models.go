package category

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoryType string

const (
	CategoryTypeIncome  CategoryType = "income"
	CategoryTypeExpense CategoryType = "expense"
)

type Category struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Name      string             `bson:"name" json:"name"`
	Type      CategoryType       `bson:"type" json:"type"`
	Icon      string             `bson:"icon" json:"icon"`
	Color     string             `bson:"color" json:"color"`
	IsDefault bool               `bson:"is_default" json:"is_default"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
