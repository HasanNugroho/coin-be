package allocation

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Allocation represents a salary allocation rule
type Allocation struct {
	ID             primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	UserID         primitive.ObjectID  `bson:"user_id" json:"user_id"`
	PocketID       *primitive.ObjectID `bson:"pocket_id,omitempty" json:"pocket_id,omitempty"`
	UserPlatformID *primitive.ObjectID `bson:"user_platform_id,omitempty" json:"user_platform_id,omitempty"`
	Priority       int                 `bson:"priority" json:"priority"` // 1=HIGH, 2=MEDIUM, 3=LOW
	AllocationType string              `bson:"allocation_type" json:"allocation_type" enums:"PERCENTAGE,NOMINAL"`
	Nominal        float64             `bson:"nominal" json:"nominal"` // percentage or amount
	IsActive       bool                `bson:"is_active" json:"is_active"`
	CreatedAt      time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time           `bson:"updated_at" json:"updated_at"`
	DeletedAt      *time.Time          `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

type AllocationType string

const (
	TypePercentage AllocationType = "PERCENTAGE"
	TypeNominal    AllocationType = "NOMINAL"
)

type Priority int

const (
	PriorityHigh   Priority = 1
	PriorityMedium Priority = 2
	PriorityLow    Priority = 3
)

func IsValidAllocationType(t string) bool {
	switch t {
	case string(TypePercentage), string(TypeNominal):
		return true
	default:
		return false
	}
}

func IsValidPriority(p int) bool {
	return p >= 1 && p <= 3
}
