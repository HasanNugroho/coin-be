package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TargetResponse struct {
	ID            primitive.ObjectID `json:"id"`
	UserID        primitive.ObjectID `json:"user_id"`
	AllocationID  primitive.ObjectID `json:"allocation_id"`
	Name          string             `json:"name"`
	TargetAmount  float64            `json:"target_amount"`
	CurrentAmount float64            `json:"current_amount"`
	Progress      float64            `json:"progress"`
	Deadline      time.Time          `json:"deadline"`
	Status        string             `json:"status"`
	CreatedAt     time.Time          `json:"created_at"`
}
