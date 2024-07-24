package helper

import "github.com/google/uuid"

type ValidateRefreshTokenParams struct {
	RefreshToken string
	UserID       uuid.UUID
	Scope        string
}

type DeActiveTokenParams struct {
	Token string
	Scope string
}

type ActiveTokenParams struct {
	Token    string
	Scope    string
	UserID   uuid.UUID
	TokenIAT int
}
