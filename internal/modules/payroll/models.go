package payroll

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PayrollRecord tracks executed payroll to ensure idempotency
type PayrollRecord struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Year      int                `bson:"year" json:"year"`
	Month     int                `bson:"month" json:"month"`
	Day       int                `bson:"day" json:"day"`
	Amount    float64            `bson:"amount" json:"amount"`
	Status    string             `bson:"status" json:"status"` // SUCCESS, FAILED
	Error     *string            `bson:"error,omitempty" json:"error,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

// PayrollStatus constants
const (
	StatusSuccess = "SUCCESS"
	StatusFailed  = "FAILED"
)
