package repository

import (
	"context"
	"oauth-server/app/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, tx *gorm.DB, user *entity.User) error
	UpdateUser(ctx context.Context, tx *gorm.DB, user *entity.User) error
	DeleteUser(ctx context.Context, tx *gorm.DB, user *entity.User) error
	BulkCreateUser(ctx context.Context, tx *gorm.DB, users []entity.User) error
	FindUserByFilter(ctx context.Context, tx *gorm.DB, filter *FindUserByFilter) (*entity.User, error)
	FindUsersByFilter(ctx context.Context, tx *gorm.DB, filer *FindUserByFilter) ([]entity.User, error)
}

type OAuthRepository interface {
	CreateOAuth(ctx context.Context, tx *gorm.DB, oauth *entity.Oauth) error
	UpdateOAuth(ctx context.Context, tx *gorm.DB, oauth *entity.Oauth) error
	FindOAuthByFilter(ctx context.Context, tx *gorm.DB, filter *FindOAuthByFilter) (*entity.Oauth, error)
}
