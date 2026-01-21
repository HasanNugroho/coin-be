package target

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TargetStatus string

const (
	TargetStatusActive    TargetStatus = "active"
	TargetStatusCompleted TargetStatus = "completed"
)

type SavingTarget struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
	AllocationID  primitive.ObjectID `bson:"allocation_id" json:"allocation_id"`
	Name          string             `bson:"name" json:"name"`
	TargetAmount  float64            `bson:"target_amount" json:"target_amount"`
	CurrentAmount float64            `bson:"current_amount" json:"current_amount"`
	Deadline      time.Time          `bson:"deadline" json:"deadline"`
	Status        TargetStatus       `bson:"status" json:"status"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
}
