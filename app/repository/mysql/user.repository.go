package mysql_repository

import (
	"gin-boilerplate/app/entity"
	"gin-boilerplate/app/repository"
	"time"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewMysqlUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{
		db,
	}
}

func (r *userRepository) NewUserTransaction() *gorm.DB {
	return r.db.Begin()
}

func (r *userRepository) CreateUser(user *entity.User) error {
	return r.db.Create(&user).Error
}

func (r *userRepository) UpdateUser(user *entity.User) error {
	user.UpdatedAt = time.Now().Unix()

	return r.db.Save(&user).Error
}

func (r *userRepository) DeleteUser(user *entity.User) error {
	return r.db.Delete(&user).Error
}

func (r *userRepository) NewUser() *entity.User {
	return &entity.User{
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
}

func (r *userRepository) BulkCreateUser(users []entity.User) error {
	return r.db.Create(&users).Error
}

func (r *userRepository) FindUserByFilter(filter *repository.FindUserByFilter) (*entity.User, error) {
	var user *entity.User
	err := r.db.First(&user).Error
	return user, err
}

func (r *userRepository) FindUsersByFilter(filer *repository.FindUserByFilter) ([]entity.User, error) {
	var user []entity.User
	err := r.db.Find(&user).Error
	return user, err
}
