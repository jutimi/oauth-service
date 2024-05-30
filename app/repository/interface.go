package repository

import (
	"gin-boilerplate/app/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	NewUserTransaction() *gorm.DB
	CreateUser(user *entity.User) error
	UpdateUser(user *entity.User) error
	DeleteUser(user *entity.User) error
	NewUser() *entity.User
	BulkCreateUser(users []entity.User) error
	FindUserByFilter(filter *FindUserByFilter) (*entity.User, error)
	FindUsersByFilter(filer *FindUserByFilter) ([]entity.User, error)
}

type OAuthRepository interface {
	NewOAuthTransaction() *gorm.DB
	CreateOAuth(oauth *entity.Oauth) error
	UpdateOAuth(oauth *entity.Oauth) error
	NewOAuth() *entity.Oauth
	FindOAuthByFilter(filter *FindOAuthByFilter) (*entity.Oauth, error)
}
