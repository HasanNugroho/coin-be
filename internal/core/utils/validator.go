package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateRequest(req any) error {
	err := validate.Struct(req)
	if err != nil {
		return formatValidationError(err)
	}
	return nil
}

func formatValidationError(err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var messages []string
		for _, fieldError := range validationErrors {
			messages = append(messages, formatFieldError(fieldError))
		}
		return fmt.Errorf("validation failed: %s", strings.Join(messages, "; "))
	}
	return err
}

func formatFieldError(fieldError validator.FieldError) string {
	field := fieldError.Field()
	tag := fieldError.Tag()
	param := fieldError.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s", field, param)
	case "max":
		return fmt.Sprintf("%s must be at most %s", field, param)
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters", field, param)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, param)
	case "hexadecimal":
		return fmt.Sprintf("%s must be a valid hexadecimal string", field)
	case "omitempty":
		return fmt.Sprintf("%s is invalid", field)
	default:
		return fmt.Sprintf("%s failed validation: %s", field, tag)
	}
}
