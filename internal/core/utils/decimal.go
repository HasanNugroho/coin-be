package utils

import (
	"math/big"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Decimal128ToFloat64 converts MongoDB Decimal128 to float64 safely
func Decimal128ToFloat64(d primitive.Decimal128) float64 {
	s := d.String()

	// strconv.ParseFloat handles scientific notation (e.g. "1E+5") natively
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		// Fallback: parse via big.Float for edge cases like "NaN", "Inf", dll
		bf, _, err := big.ParseFloat(s, 10, 64, big.ToNearestEven)
		if err != nil {
			return 0
		}
		result, _ := bf.Float64()
		return result
	}
	return f
}

// NewDecimal128FromFloat converts float64 to MongoDB Decimal128
func NewDecimal128FromFloat(v float64) primitive.Decimal128 {
	// FormatFloat with 'f' (no scientific notation) dan presisi 2 desimal
	d, err := primitive.ParseDecimal128(strconv.FormatFloat(v, 'f', 2, 64))
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
