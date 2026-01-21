package allocation

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Allocation struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
	Name          string             `bson:"name" json:"name"`
	Priority      int                `bson:"priority" json:"priority"`
	Percentage    float64            `bson:"percentage" json:"percentage"`
	CurrentAmount float64            `bson:"current_amount" json:"current_amount"`
	TargetAmount  *float64           `bson:"target_amount,omitempty" json:"target_amount,omitempty"`
	IsActive      bool               `bson:"is_active" json:"is_active"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
}

type AllocationLog struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID          primitive.ObjectID `bson:"user_id" json:"user_id"`
	AllocationID    primitive.ObjectID `bson:"allocation_id" json:"allocation_id"`
	TransactionID   primitive.ObjectID `bson:"transaction_id" json:"transaction_id"`
	IncomeAmount    float64            `bson:"income_amount" json:"income_amount"`
	AllocatedAmount float64            `bson:"allocated_amount" json:"allocated_amount"`
	Percentage      float64            `bson:"percentage" json:"percentage"`
	Priority        int                `bson:"priority" json:"priority"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
}
