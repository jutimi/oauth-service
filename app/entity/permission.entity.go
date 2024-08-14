package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model
	Id              uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	WorkspaceId     *uuid.UUID `json:"workspace_id" gorm:"type:uuid"`
	UserWorkspaceId *uuid.UUID `json:"user_workspace_id" gorm:"type:uuid"`
	UserId          uuid.UUID  `json:"user_id" gorm:"type:uuid;not null"`
	Scopes          string     `json:"scopes" gorm:"type:text;not null"`
	CreatedAt       int64      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       int64      `json:"updated_at" gorm:"autoUpdateTime:milli"`
}

func NewPermission() *Permission {
	return &Permission{
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
}

func (u *Permission) TableName() string {
	return "permissions"
}
