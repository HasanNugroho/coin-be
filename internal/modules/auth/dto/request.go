package dto

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=255"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=255"`
	Name     string `json:"name" validate:"required,min=1,max=255"`
	Phone    string `json:"phone" validate:"omitempty,max=20"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required,min=10"`
}
