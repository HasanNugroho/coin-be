package constants

const (
	RoleAdmin   = "admin"
	RoleUser    = "user"
	RolePremium = "premium"
)

const (
	SalaryCycleMonthly   = "monthly"
	SalaryCycleWeekly    = "weekly"
	SalaryCycleBiweekly  = "biweekly"
)

const (
	DefaultLimit = 10
	MaxLimit     = 100
)

const (
	ErrorCodeInvalidInput      = "INVALID_INPUT"
	ErrorCodeUnauthorized      = "UNAUTHORIZED"
	ErrorCodeForbidden         = "FORBIDDEN"
	ErrorCodeNotFound          = "NOT_FOUND"
	ErrorCodeConflict          = "CONFLICT"
	ErrorCodeInternalServer    = "INTERNAL_SERVER_ERROR"
)
