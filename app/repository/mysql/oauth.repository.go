package mysql_repository

import (
	"context"
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

func (r *oAuthRepository) CreateOAuth(
	ctx context.Context,
	tx *gorm.DB,
	oauth *entity.Oauth,
) error {
	return r.db.Create(&oauth).Error
}

func (r *oAuthRepository) UpdateOAuth(
	ctx context.Context,
	tx *gorm.DB,
	oauth *entity.Oauth,
) error {
	oauth.UpdatedAt = time.Now().Unix()

	return r.db.Save(&oauth).Error
}

func (r *oAuthRepository) FindOAuthByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindOAuthByFilter,
) (*entity.Oauth, error) {
	var data *entity.Oauth
	err := r.db.First(&data).Error
	return data, err
}
