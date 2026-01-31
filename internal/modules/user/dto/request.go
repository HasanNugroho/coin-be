package dto

type CreateUserRequest struct {
	Email string `json:"email" validate:"required,email"`
	Phone string `json:"phone" validate:"omitempty,max=20"`
	Name  string `json:"name" validate:"required,min=1,max=255"`
}

type UpdateUserRequest struct {
	Name        string  `json:"name" validate:"omitempty,min=1,max=255"`
	Email       string  `json:"email" validate:"omitempty,email"`
	Phone       string  `json:"phone" validate:"omitempty,max=20"`
	TelegramId  string  `json:"telegramId" validate:"omitempty,max=100"`
	Currency    string  `json:"currency" validate:"omitempty,len=3"`
	BaseSalary  float64 `json:"baseSalary" validate:"omitempty,min=0"`
	SalaryCycle string  `json:"salaryCycle" validate:"omitempty,oneof=monthly weekly biweekly"`
	SalaryDay   int     `json:"salaryDay" validate:"omitempty,min=1,max=28"`
	Language    string  `json:"language" validate:"omitempty,len=2"`
}

type CreateUserProfileRequest struct {
	BaseSalary  float64 `json:"base_salary" validate:"required,min=0"`
	SalaryCycle string  `json:"salary_cycle" validate:"required,oneof=monthly weekly biweekly"`
	SalaryDay   int     `json:"salary_day" validate:"required,min=1,max=28"`
	PayCurrency string  `json:"pay_currency" validate:"omitempty,len=3"`
}

type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"omitempty,max=500"`
}

type AssignRoleRequest struct {
	RoleID string `json:"role_id" validate:"required,len=24,hexadecimal"`
}
