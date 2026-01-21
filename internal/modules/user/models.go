package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email        string             `bson:"email" json:"email"`
	Phone        string             `bson:"phone" json:"phone"`
	PasswordHash string             `bson:"password_hash" json:"-"`
	Name         string             `bson:"name" json:"name"`
	Role         string             `bson:"role" json:"role" enums:"admin,user"` // Role: admin or user
	IsActive     bool               `bson:"is_active" json:"is_active"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

type FinancialProfile struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	BaseSalary  float64            `bson:"base_salary" json:"base_salary"`
	SalaryCycle string             `bson:"salary_cycle" json:"salary_cycle"`
	SalaryDay   int                `bson:"salary_day" json:"salary_day"`
	PayCurrency string             `bson:"pay_currency" json:"pay_currency"`
	IsActive    bool               `bson:"is_active" json:"is_active"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// Role constants
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)
