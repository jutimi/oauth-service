package postgres_repository

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
		return tx.WithContext(ctx).Create(&oauth).Error
	}

	return r.db.WithContext(ctx).Create(&oauth).Error
}

func (r *oAuthRepository) UpdateOAuth(
	ctx context.Context,
	tx *gorm.DB,
	oauth *entity.Oauth,
) error {
	oauth.UpdatedAt = time.Now().Unix()

	if tx != nil {
		return tx.WithContext(ctx).Save(&oauth).Error
	}

	return r.db.WithContext(ctx).Save(&oauth).Error
}

func (r *oAuthRepository) FindOAuthByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindOAuthByFilter,
) (*entity.Oauth, error) {
	var data *entity.Oauth

	query := r.db.WithContext(ctx).Debug()
	if tx != nil {
		query = tx.WithContext(ctx).Debug()
	}

	if filter.Token != nil {
		query = query.Scopes(findByText(*filter.Token, "token"))
	}
	if filter.UserID != nil {
		query = query.Scopes(findById(*filter.UserID, "user_id"))
	}
	if filter.PlatForm != nil {
		query = query.Scopes(findByText(*filter.PlatForm, "platform"))
	}

	err := query.First(&data).Error
	return data, err
}
