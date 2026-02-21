package dto

import "time"

type UserResponse struct {
	ID                       string    `json:"id"`
	Name                     string    `json:"name"`
	Email                    string    `json:"email"`
	Phone                    string    `json:"phone"`
	TelegramId               string    `json:"telegramId"`
	Currency                 string    `json:"currency"`
	BaseSalary               float64   `json:"baseSalary"`
	SalaryCycle              string    `json:"salaryCycle"`
	SalaryDay                int       `json:"salaryDay"`
	Language                 string    `json:"language"`
	AutoInputPayroll         bool      `json:"autoInputPayroll"`
	DefaultUserPlatformID    *string   `json:"defaultUserPlatformId,omitempty"`
	TelegramIntegrationAlert bool      `json:"telegramIntegrationAlert"`
	IsActive                 bool      `json:"is_active"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
}

type UserProfileResponse struct {
	ID                    string    `json:"id"`
	UserID                string    `json:"user_id"`
	BaseSalary            float64   `json:"base_salary"`
	SalaryCycle           string    `json:"salary_cycle"`
	SalaryDay             int       `json:"salary_day"`
	PayCurrency           string    `json:"pay_currency"`
	AutoInputPayroll      bool      `json:"auto_input_payroll"`
	DefaultUserPlatformID *string   `json:"default_user_platform_id,omitempty"`
	IsActive              bool      `json:"is_active"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type RoleResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

type UserRoleResponse struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	RoleID     string    `json:"role_id"`
	AssignedAt time.Time `json:"assigned_at"`
}

type ListUsersResponse struct {
	Users []UserResponse `json:"users"`
	Total int64          `json:"total"`
	Limit int64          `json:"limit"`
	Skip  int64          `json:"skip"`
}
