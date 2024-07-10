package postgres_repository

import (
	"context"
	"oauth-server/app/entity"
	"oauth-server/app/repository"
	"time"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{
		db,
	}
}

func (r *userRepository) Create(
	ctx context.Context,
	tx *gorm.DB,
	user *entity.User,
) error {
	if tx != nil {
		return tx.WithContext(ctx).Create(&user).Error
	}

	return r.db.WithContext(ctx).Create(&user).Error
}

func (r *userRepository) Update(
	ctx context.Context,
	tx *gorm.DB,
	user *entity.User,
) error {
	user.UpdatedAt = time.Now().Unix()

	if tx != nil {
		return tx.WithContext(ctx).Save(&user).Error
	}

	return r.db.WithContext(ctx).Save(&user).Error
}

func (r *userRepository) Delete(
	ctx context.Context,
	tx *gorm.DB,
	user *entity.User,
) error {
	if tx != nil {
		return tx.WithContext(ctx).Delete(&user).Error
	}

	return r.db.WithContext(ctx).Delete(&user).Error
}

func (r *userRepository) BulkCreate(
	ctx context.Context,
	tx *gorm.DB,
	users []entity.User,
) error {
	if tx != nil {
		return tx.WithContext(ctx).Create(&users).Error
	}

	return r.db.WithContext(ctx).Create(&users).Error
}

func (r *userRepository) FindOneByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindUserByFilter,
) (*entity.User, error) {
	var data *entity.User
	query := r.buildFilter(ctx, tx, filter)

	err := query.First(&data).Error
	return data, err
}

func (r *userRepository) FindByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filer *repository.FindUserByFilter,
) ([]entity.User, error) {
	var data []entity.User
	query := r.buildFilter(ctx, tx, filer)

	err := query.Find(&data).Error
	return data, err
}

func (r *userRepository) buildFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindUserByFilter,
) *gorm.DB {
	query := r.db.WithContext(ctx)
	if tx != nil {
		query = tx.WithContext(ctx)
	}

	if filter.Email != nil && *filter.Email != "" {
		query = query.Scopes(findByText(*filter.Email, "email"))
	}
	if filter.PhoneNumber != nil && *filter.PhoneNumber != "" {
		query = query.Scopes(findByText(*filter.PhoneNumber, "phone_number"))
	}
	if filter.ID != nil {
		query = query.Scopes(findById(*filter.ID, "id"))
	}
	if filter.IDs != nil && len(filter.IDs) > 0 {
		query = query.Scopes(findBySlice(filter.IDs, "id"))
	}
	if filter.Emails != nil && len(filter.Emails) > 0 {
		query = query.Scopes(findBySlice(filter.Emails, "email"))
	}
	if filter.PhoneNumbers != nil && len(filter.PhoneNumbers) > 0 {
		query = query.Scopes(findBySlice(filter.PhoneNumbers, "phone_number"))
	}
	if filter.Limit != nil && filter.Offset != nil {
		query = query.Scopes(paginate(*filter.Limit, *filter.Offset))
	}

	return query
}
