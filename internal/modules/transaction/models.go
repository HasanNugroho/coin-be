package transaction

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID         primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID  `bson:"user_id" json:"user_id"`
	Type       string              `bson:"type" json:"type" enums:"income,expense,transfer"`
	Amount     float64             `bson:"amount" json:"amount"`
	PocketFrom *primitive.ObjectID `bson:"pocket_from,omitempty" json:"pocket_from,omitempty"`
	PocketTo   *primitive.ObjectID `bson:"pocket_to,omitempty" json:"pocket_to,omitempty"`
	// UserPlatformFrom and UserPlatformTo hold user-specific platform references for balance updates
	UserPlatformFrom *primitive.ObjectID `bson:"user_platform_from,omitempty" json:"user_platform_from,omitempty"`
	UserPlatformTo   *primitive.ObjectID `bson:"user_platform_to,omitempty" json:"user_platform_to,omitempty"`
	CategoryID       *primitive.ObjectID `bson:"category_id,omitempty" json:"category_id,omitempty"`
	// PlatformID is reference-only (AdminPlatform), not used for balance logic
	PlatformID *primitive.ObjectID `bson:"platform_id,omitempty" json:"platform_id,omitempty"`
	Note       *string             `bson:"note,omitempty" json:"note,omitempty"`
	Date       time.Time           `bson:"date" json:"date"`
	Ref        *string             `bson:"ref,omitempty" json:"ref,omitempty"`

	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

type TransactionType string

const (
	TypeIncome   TransactionType = "income"
	TypeExpense  TransactionType = "expense"
	TypeTransfer TransactionType = "transfer"
)

func IsValidTransactionType(t string) bool {
	switch t {
	case string(TypeIncome), string(TypeExpense), string(TypeTransfer):
		return true
	default:
		return false
	}
}
