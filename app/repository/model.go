package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FindUserByFilter struct {
	Email        *string
	PhoneNumber  *string
	ID           *uuid.UUID
	IDs          []uuid.UUID
	Emails       []string
	PhoneNumbers []string
	Limit        *int
	Offset       *int
	IsActive     *bool
}

type FindOAuthByFilter struct {
	UserID   *uuid.UUID
	Token    *string
	PlatForm *string
	Scope    *string
}

type FindByFilterForUpdateParams struct {
	Filter     interface{}
	LockOption string
	Tx         *gorm.DB
}
