package entity

import "github.com/google/uuid"

const (
	OAuthPlatformMobile = "mobile"
	OAuthPlatformWeb    = "web"
)

const (
	OAuthStatusActive   = "active"
	OAuthStatusInactive = "inactive"
	OAuthStatusBlocked  = "blocked"
)

type Oauth struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserId    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	IP        string    `json:"ip" gorm:"type:text;not null"`
	Platform  string    `json:"platform" gorm:"type:varchar(10);not null"`
	Token     string    `json:"token" gorm:"type:text;not null"`
	Status    bool      `json:"status" gorm:"varchar(10);not null"`
	CreatedAt int64     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64     `json:"updated_at" gorm:"autoUpdateTime:milli"`
}

func (Oauth) TableName() string {
	return "oauth"
}
