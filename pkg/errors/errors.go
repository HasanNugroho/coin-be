package errors

type ErrorResponse struct {
	Error string `json:"error"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrorResponse struct {
	Errors []ValidationError `json:"errors"`
}

func NewErrorResponse(message string) ErrorResponse {
	return ErrorResponse{Error: message}
}

func NewValidationErrorResponse(errors []ValidationError) ValidationErrorResponse {
	return ValidationErrorResponse{Errors: errors}
}

type AppError struct {
	Code    string
	Message string
	Status  int
}

func NewAppError(code, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

func (e *AppError) Error() string {
	return e.Message
}
