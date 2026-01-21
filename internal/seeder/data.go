package seeder

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id,omitempty"`
	Name      string             `bson:"name"`
	Type      string             `bson:"type"` // "income" or "expense"
	Icon      string             `bson:"icon"`
	Color     string             `bson:"color"`
	IsDefault bool               `bson:"is_default"`
	CreatedAt time.Time          `bson:"created_at"`
}

type Allocation struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	UserID        primitive.ObjectID `bson:"user_id,omitempty"`
	Name          string             `bson:"name"`
	Priority      int                `bson:"priority"`
	Percentage    float64            `bson:"percentage"`
	CurrentAmount float64            `bson:"current_amount"`
	TargetAmount  float64            `bson:"target_amount,omitempty"`
	IsActive      bool               `bson:"is_active"`
	CreatedAt     time.Time          `bson:"created_at"`
}

func getDefaultCategories() []Category {
	now := time.Now()

	return []Category{
		// Income categories
		{
			Name:      "Salary",
			Type:      "income",
			Icon:      "💼",
			Color:     "#3498db",
			IsDefault: true,
			CreatedAt: now,
		},
		{
			Name:      "Bonus",
			Type:      "income",
			Icon:      "🎁",
			Color:     "#2ecc71",
			IsDefault: true,
			CreatedAt: now,
		},
		{
			Name:      "Investment",
			Type:      "income",
			Icon:      "📈",
			Color:     "#f39c12",
			IsDefault: true,
			CreatedAt: now,
		},
		{
			Name:      "Freelance",
			Type:      "income",
			Icon:      "💻",
			Color:     "#9b59b6",
			IsDefault: true,
			CreatedAt: now,
		},
		// Expense categories
		{
			Name:      "Food & Dining",
			Type:      "expense",
			Icon:      "🍔",
			Color:     "#e74c3c",
			IsDefault: true,
			CreatedAt: now,
		},
		{
			Name:      "Transportation",
			Type:      "expense",
			Icon:      "🚗",
			Color:     "#34495e",
			IsDefault: true,
			CreatedAt: now,
		},
		{
			Name:      "Shopping",
			Type:      "expense",
			Icon:      "🛍️",
			Color:     "#e91e63",
			IsDefault: true,
			CreatedAt: now,
		},
		{
			Name:      "Bills & Utilities",
			Type:      "expense",
			Icon:      "💡",
			Color:     "#1abc9c",
			IsDefault: true,
			CreatedAt: now,
		},
		{
			Name:      "Entertainment",
			Type:      "expense",
			Icon:      "🎬",
			Color:     "#c0392b",
			IsDefault: true,
			CreatedAt: now,
		},
		{
			Name:      "Healthcare",
			Type:      "expense",
			Icon:      "🏥",
			Color:     "#16a085",
			IsDefault: true,
			CreatedAt: now,
		},
		{
			Name:      "Education",
			Type:      "expense",
			Icon:      "📚",
			Color:     "#2980b9",
			IsDefault: true,
			CreatedAt: now,
		},
		{
			Name:      "Personal Care",
			Type:      "expense",
			Icon:      "💅",
			Color:     "#d35400",
			IsDefault: true,
			CreatedAt: now,
		},
	}
}

func getDefaultAllocations() []Allocation {
	now := time.Now()

	return []Allocation{
		{
			Name:          "Bills & Utilities",
			Priority:      1,
			Percentage:    40,
			CurrentAmount: 0,
			IsActive:      true,
			CreatedAt:     now,
		},
		{
			Name:          "Emergency Fund",
			Priority:      2,
			Percentage:    10,
			CurrentAmount: 0,
			TargetAmount:  10000000,
			IsActive:      true,
			CreatedAt:     now,
		},
		{
			Name:          "Investment",
			Priority:      3,
			Percentage:    30,
			CurrentAmount: 0,
			IsActive:      true,
			CreatedAt:     now,
		},
		{
			Name:          "Savings",
			Priority:      4,
			Percentage:    20,
			CurrentAmount: 0,
			TargetAmount:  5000000,
			IsActive:      true,
			CreatedAt:     now,
		},
	}
}
