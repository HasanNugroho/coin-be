package daily_summary

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DailySummary struct {
	ID                primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	UserID            primitive.ObjectID  `bson:"user_id" json:"user_id"`
	Date              time.Time           `bson:"date" json:"date"`
	TotalIncome       float64             `bson:"total_income" json:"total_income"`
	TotalExpense      float64             `bson:"total_expense" json:"total_expense"`
	CategoryBreakdown []CategoryBreakdown `bson:"category_breakdown" json:"category_breakdown"`
	PocketBreakdown   []PocketBreakdown   `bson:"pocket_breakdown" json:"pocket_breakdown"`
	PlatformBreakdown []PlatformBreakdown `bson:"platform_breakdown" json:"platform_breakdown"`
	CreatedAt         time.Time           `bson:"created_at" json:"created_at"`
}

type CategoryBreakdown struct {
	CategoryID   *primitive.ObjectID `bson:"category_id,omitempty" json:"category_id,omitempty"`
	CategoryName string              `bson:"category_name" json:"category_name"`
	Type         string              `bson:"type" json:"type"`
	Amount       float64             `bson:"amount" json:"amount"`
}

type PocketBreakdown struct {
	PocketID   *primitive.ObjectID `bson:"pocket_id,omitempty" json:"pocket_id,omitempty"`
	PocketName string              `bson:"pocket_name" json:"pocket_name"`
	Type       string              `bson:"type" json:"type"`
	Amount     float64             `bson:"amount" json:"amount"`
}

type PlatformBreakdown struct {
	PlatformID   *primitive.ObjectID `bson:"platform_id,omitempty" json:"platform_id,omitempty"`
	PlatformName string              `bson:"platform_name" json:"platform_name"`
	Type         string              `bson:"type" json:"type"`
	Amount       float64             `bson:"amount" json:"amount"`
}
