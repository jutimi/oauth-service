package mysql_repository

import (
	"context"
	"oauth-server/app/entity"
	"oauth-server/app/repository"
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
	return r.db.WithContext(ctx).Create(&oauth).Error
}

func (r *oAuthRepository) UpdateOAuth(
	ctx context.Context,
	tx *gorm.DB,
	oauth *entity.Oauth,
) error {
	oauth.UpdatedAt = time.Now().Unix()

	return r.db.WithContext(ctx).Save(&oauth).Error
}

func (r *oAuthRepository) FindOAuthByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindOAuthByFilter,
) (*entity.Oauth, error) {
	var data *entity.Oauth
	err := r.db.WithContext(ctx).First(&data).Error
	return data, err
}
