package mysql_repository

import (
	"gin-boilerplate/app/entity"
	"gin-boilerplate/app/repository"
	"time"

	"gorm.io/gorm"
)

type oAuthRepository struct {
	db *gorm.DB
}

func NewMysqlOAuthRepository(db *gorm.DB) repository.OAuthRepository {
	return &oAuthRepository{
		db,
	}
}

func (r *oAuthRepository) NewOAuthTransaction() *gorm.DB {
	return r.db.Begin()
}

func (r *oAuthRepository) CreateOAuth(oauth *entity.Oauth) error {
	return r.db.Create(&oauth).Error
}

func (r *oAuthRepository) UpdateOAuth(oauth *entity.Oauth) error {
	oauth.UpdatedAt = time.Now().Unix()

	return r.db.Save(&oauth).Error
}
func (r *oAuthRepository) NewOAuth() *entity.Oauth {
	return &entity.Oauth{
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
}

func (r *oAuthRepository) FindOAuthByFilter(filter *repository.FindOAuthByFilter) (*entity.Oauth, error) {
	var data *entity.Oauth
	err := r.db.First(&data).Error
	return data, err
}
