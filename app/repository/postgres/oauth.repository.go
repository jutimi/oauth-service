package postgres_repository

import (
	"context"
	"errors"
	"oauth-server/app/entity"
	"oauth-server/app/repository"
	"time"

	"gorm.io/gorm"
)

type oAuthRepository struct {
	db *gorm.DB
}

func NewOAuthRepository(db *gorm.DB) repository.OAuthRepository {
	return &oAuthRepository{
		db,
	}
}

func (r *oAuthRepository) Create(
	ctx context.Context,
	tx *gorm.DB,
	oauth *entity.Oauth,
) error {
	if tx != nil {
		return tx.WithContext(ctx).Create(&oauth).Error
	}

	return r.db.WithContext(ctx).Create(&oauth).Error
}

func (r *oAuthRepository) Update(
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

func (r *oAuthRepository) FindOneByFilter(
	ctx context.Context,
	filter *repository.FindOAuthByFilter,
) (*entity.Oauth, error) {
	var data *entity.Oauth
	query := r.buildFilter(ctx, nil, filter)

	err := query.First(&data).Error
	return data, err
}

func (r *oAuthRepository) FindOneByFilterForUpdate(
	ctx context.Context,
	data *repository.FindByFilterForUpdateParams,
) (*entity.Oauth, error) {
	filter, ok := data.Filter.(*repository.FindOAuthByFilter)
	if !ok {
		return nil, errors.New("invalid argument")
	}

	var oauth *entity.Oauth
	query := r.buildFilter(ctx, data.Tx, filter)
	query = buildLockQuery(query, data.LockOption)

	err := query.First(&oauth).Error
	return oauth, err
}

// -------------------------------------------------------------------------------
func (r *oAuthRepository) buildFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindOAuthByFilter,
) *gorm.DB {
	query := r.db.WithContext(ctx)
	if tx != nil {
		query = tx.WithContext(ctx)
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

	return query
}
