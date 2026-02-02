package user_platform

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserPlatform represents a user-owned platform with real-time balance.
// This is the entity that holds platform balance for a specific user.
// Platform is reference-only and does not hold balance.
type UserPlatform struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	PlatformID primitive.ObjectID `bson:"platform_id" json:"platform_id"`

	// AliasName allows user-friendly naming (e.g., "BRI - Salary", "BRI - Saving")
	// If null, fallback to Platform.name on response layer
	AliasName *string `bson:"alias_name,omitempty" json:"alias_name,omitempty"`

	// Balance is user-specific and updated through transactions
	Balance primitive.Decimal128 `bson:"balance" json:"balance"`

	IsActive  bool       `bson:"is_active" json:"is_active"`
	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}
