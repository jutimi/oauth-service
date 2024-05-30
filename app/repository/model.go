package repository

import "github.com/google/uuid"

type FindUserByFilter struct {
	Email       string
	PhoneNumber string
}

type FindOAuthByFilter struct {
	UserId   uuid.UUID
	Token    string
	PlatForm string
}
