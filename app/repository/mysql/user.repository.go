package mysql_repository

import (
	"context"
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

func (r *userRepository) CreateUser(
	ctx context.Context,
	tx *gorm.DB,
	user *entity.User,
) error {
	return r.db.Create(&user).Error
}

func (r *userRepository) UpdateUser(
	ctx context.Context,
	tx *gorm.DB,
	user *entity.User,
) error {
	user.UpdatedAt = time.Now().Unix()

	return r.db.Save(&user).Error
}

func (r *userRepository) DeleteUser(
	ctx context.Context,
	tx *gorm.DB,
	user *entity.User,
) error {
	return r.db.Delete(&user).Error
}

func (r *userRepository) BulkCreateUser(
	ctx context.Context,
	tx *gorm.DB,
	users []entity.User,
) error {
	return r.db.Create(&users).Error
}

func (r *userRepository) FindUserByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindUserByFilter,
) (*entity.User, error) {
	var user *entity.User
	err := r.db.First(&user).Error
	return user, err
}

func (r *userRepository) FindUsersByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filer *repository.FindUserByFilter,
) ([]entity.User, error) {
	var user []entity.User
	err := r.db.Find(&user).Error
	return user, err
}
