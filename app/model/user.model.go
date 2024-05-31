package model

type LoginRequest struct {
	Password    string `json:"password" binding:"required"`
	Email       string `json:"email" binding:"omitempty,email"`
	PhoneNumber string `json:"phone_number" binding:"omitempty,phone_number"`
}
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterRequest struct {
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=password"`
	Email           string `json:"email" binding:"omitempty,email"`
	PhoneNumber     string `json:"phone_number" binding:"omitempty,phone_number"`
}
type RegisterResponse struct {
}

type LogoutRequest struct {
}
type LogoutResponse struct {
}
