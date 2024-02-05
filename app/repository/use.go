package repository

import (
	"gin-boilerplate/app/entity"
	"time"

	"gorm.io/gorm"
)

type FindUserByFilter struct {
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db,
	}
}

func (r *userRepository) Create(user *entity.User) error {
	return r.db.Create(&user).Error
}

func (r *userRepository) Update(user *entity.User) error {
	user.UpdatedAt = time.Now().Unix()

	return r.db.Save(&user).Error
}

func (r *userRepository) Delete(user *entity.User) error {
	return r.db.Delete(&user).Error
}

func (r *userRepository) New() *entity.User {
	return &entity.User{
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
}

func (r *userRepository) BulkCreate(users []entity.User) error {
	return r.db.Create(&users).Error
}

func (r *userRepository) FindOneByFilter(filter *FindUserByFilter) (*entity.User, error) {
	var user *entity.User
	err := r.db.First(&user).Error
	return user, err
}
