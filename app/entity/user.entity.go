package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID          uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Email       *string   `json:"email" gorm:"type:varchar(100);"`
	PhoneNumber *string   `json:"phone_number" gorm:"type:varchar(20);"`
	Password    string    `json:"password" gorm:"type:text;not null"`
	IsActive    bool      `json:"is_active" gorm:"default:true;type:bool;not null"`
	CreatedAt   int64     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   int64     `json:"updated_at" gorm:"autoUpdateTime:milli"`
}

func (User) TableName() string {
	return "users"
}
