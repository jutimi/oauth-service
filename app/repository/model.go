package repository

import "github.com/google/uuid"

type FindUserByFilter struct {
	Email       *string
	PhoneNumber *string
	ID          *uuid.UUID
}

type FindOAuthByFilter struct {
	UserID   *uuid.UUID
	Token    *string
	PlatForm *string
}
