package model

type LoginRequest struct {
	Password    string `json:"password" binding:"required"`
	Email       string `json:"email" binding:"omitempty,email"`
	PhoneNumber string `json:"phone_number" binding:"omitempty,phone_number"`
}
type LoginResponse struct {
}

type RegisterRequest struct {
}
type RegisterResponse struct {
}

type LogoutRequest struct {
}
type LogoutResponse struct {
}
