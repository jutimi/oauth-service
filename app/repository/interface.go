package repository

import (
	"context"
	"oauth-server/app/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, tx *gorm.DB, user *entity.User) error
	Update(ctx context.Context, tx *gorm.DB, user *entity.User) error
	Delete(ctx context.Context, tx *gorm.DB, user *entity.User) error
	BulkCreate(ctx context.Context, tx *gorm.DB, users []entity.User) error
	FindOneByFilter(ctx context.Context, tx *gorm.DB, filter *FindUserByFilter) (*entity.User, error)
	FindByFilter(ctx context.Context, tx *gorm.DB, filer *FindUserByFilter) ([]entity.User, error)
	FindExistedByFilter(ctx context.Context, tx *gorm.DB, filter *FindUserByFilter) ([]entity.User, error)
}

type OAuthRepository interface {
	Create(ctx context.Context, tx *gorm.DB, oauth *entity.Oauth) error
	Update(ctx context.Context, tx *gorm.DB, oauth *entity.Oauth) error
	FindOneByFilter(ctx context.Context, tx *gorm.DB, filter *FindOAuthByFilter) (*entity.Oauth, error)
}
