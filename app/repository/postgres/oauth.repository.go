package postgres_repository

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

func NewPostgresOAuthRepository(db *gorm.DB) repository.OAuthRepository {
	return &oAuthRepository{
		db,
	}
}

func (r *oAuthRepository) CreateOAuth(
	ctx context.Context,
	tx *gorm.DB,
	oauth *entity.Oauth,
) error {
	if tx != nil {
		return tx.Create(&oauth).Error
	}

	return r.db.Create(&oauth).Error
}

func (r *oAuthRepository) UpdateOAuth(
	ctx context.Context,
	tx *gorm.DB,
	oauth *entity.Oauth,
) error {
	oauth.UpdatedAt = time.Now().Unix()

	if tx != nil {
		return tx.Save(&oauth).Error
	}

	return r.db.Save(&oauth).Error
}

func (r *oAuthRepository) FindOAuthByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindOAuthByFilter,
) (*entity.Oauth, error) {
	var data *entity.Oauth

	if tx != nil {
		err := tx.First(&data).Error
		return data, err
	}

	err := r.db.First(&data).Error
	return data, err
}
