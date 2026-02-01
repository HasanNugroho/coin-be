package utils

import (
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Decimal128ToFloat64 converts MongoDB Decimal128 to float64 safely
func Decimal128ToFloat64(d primitive.Decimal128) float64 {
	f, err := strconv.ParseFloat(d.String(), 64)
	if err != nil {
		return 0
	}
	return f
}

// NewDecimal128FromFloat converts float64 to MongoDB Decimal128
func NewDecimal128FromFloat(v float64) primitive.Decimal128 {
	d, err := primitive.ParseDecimal128(strconv.FormatFloat(v, 'f', -1, 64))
	if err != nil {
		return primitive.NewDecimal128(0, 0)
	}
	return d
}

// AddDecimal128 adds float64 amount to Decimal128
func AddDecimal128(d primitive.Decimal128, amount float64) primitive.Decimal128 {
	current := Decimal128ToFloat64(d)
	return NewDecimal128FromFloat(current + amount)
}
