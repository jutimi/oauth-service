package repository

import "gin-boilerplate/app/entity"

type UserRepository interface {
	CreateUser(user *entity.User) error
	UpdateUser(user *entity.User) error
	DeleteUser(user *entity.User) error
	NewUser() *entity.User
	BulkCreateUser(users []entity.User) error
	FindUserByFilter(filter *FindUserByFilter) (*entity.User, error)
}
