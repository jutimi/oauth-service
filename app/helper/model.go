package helper

import "github.com/google/uuid"

// Params func ValidateRefreshToken
type ValidateRefreshTokenParams struct {
	RefreshToken string
	UserId       uuid.UUID
	Scope        string
}

// Params func DeActiveToken
type DeActiveTokenParams struct {
	Token string
	Scope string
}

// Params func ActiveToken
type ActiveTokenParams struct {
	Token    string
	Scope    string
	UserId   uuid.UUID
	TokenIAT int
}

// Params func CreateUser
type CreateUserParams struct {
	PhoneNumber *string
	Email       *string
	Password    string
}
