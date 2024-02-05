package repository

import "gin-boilerplate/app/entity"

type UserRepository interface {
	Create(user *entity.User) error
	Update(user *entity.User) error
	Delete(user *entity.User) error
	New() *entity.User
	BulkCreate(users []entity.User) error
	FindOneByFilter(filter *FindUserByFilter) (*entity.User, error)
}
