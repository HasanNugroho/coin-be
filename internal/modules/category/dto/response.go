package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoryResponse struct {
	ID        primitive.ObjectID `json:"id"`
	UserID    primitive.ObjectID `json:"user_id"`
	Name      string             `json:"name"`
	Type      string             `json:"type"`
	Icon      string             `json:"icon"`
	Color     string             `json:"color"`
	IsDefault bool               `json:"is_default"`
	CreatedAt time.Time          `json:"created_at"`
}
