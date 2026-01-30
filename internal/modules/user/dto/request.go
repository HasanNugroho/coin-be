package dto

type CreateUserRequest struct {
	Email string `json:"email" binding:"required,email"`
	Phone string `json:"phone"`
	Name  string `json:"name" binding:"required"`
}

type UpdateUserRequest struct {
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	Phone       string  `json:"phone"`
	TelegramId  string  `json:"telegramId"`
	Currency    string  `json:"currency"`
	BaseSalary  float64 `json:"baseSalary"`
	SalaryCycle string  `json:"salaryCycle"`
	SalaryDay   int     `json:"salaryDay"`
	Language    string  `json:"language"`
}

type CreateUserProfileRequest struct {
	BaseSalary  float64 `json:"base_salary" binding:"required"`
	SalaryCycle string  `json:"salary_cycle" binding:"required,oneof=monthly weekly biweekly"`
	SalaryDay   int     `json:"salary_day" binding:"required,min=1,max=28"`
	PayCurrency string  `json:"pay_currency"`
}

type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type AssignRoleRequest struct {
	RoleID string `json:"role_id" binding:"required"`
}
