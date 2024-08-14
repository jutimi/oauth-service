package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FindUserByFilter struct {
	Email        *string
	PhoneNumber  *string
	Id           *uuid.UUID
	Ids          []uuid.UUID
	Emails       []string
	PhoneNumbers []string
	Limit        *int
	Offset       *int
	IsActive     *bool
}

type FindOAuthByFilter struct {
	UserId   *uuid.UUID
	Token    *string
	PlatForm *string
	Scope    *string
}

type FindByFilterForUpdateParams struct {
	Filter     interface{}
	LockOption string
	Tx         *gorm.DB
}

type FindPermissionByFilter struct {
	WorkspaceId     *uuid.UUID
	UserWorkspaceId *uuid.UUID
	UserId          *uuid.UUID
	Permission      *string
}
