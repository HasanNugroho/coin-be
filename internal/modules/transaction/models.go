package transaction

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionType string

const (
	TransactionTypeIncome  TransactionType = "income"
	TransactionTypeExpense TransactionType = "expense"
)

type Transaction struct {
	ID              primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	UserID          primitive.ObjectID  `bson:"user_id" json:"user_id"`
	Type            TransactionType     `bson:"type" json:"type"`
	Amount          float64             `bson:"amount" json:"amount"`
	CategoryID      primitive.ObjectID  `bson:"category_id" json:"category_id"`
	AllocationID    *primitive.ObjectID `bson:"allocation_id,omitempty" json:"allocation_id,omitempty"`
	Description     string              `bson:"description" json:"description"`
	TransactionDate time.Time           `bson:"transaction_date" json:"transaction_date"`
	IsDistributed   bool                `bson:"is_distributed" json:"is_distributed"`
	CreatedAt       time.Time           `bson:"created_at" json:"created_at"`
}
