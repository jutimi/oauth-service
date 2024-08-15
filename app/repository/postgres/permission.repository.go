package postgres_repository

import (
	"context"
	"oauth-server/app/entity"
	"oauth-server/app/repository"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) repository.PermissionRepository {
	return &permissionRepository{
		db,
	}
}

func (r *permissionRepository) Create(
	ctx context.Context,
	tx *gorm.DB,
	permission *entity.Permission,
) error {
	if tx != nil {
		return tx.WithContext(ctx).Create(&permission).Error
	}

	return r.db.WithContext(ctx).Create(&permission).Error
}

func (r *permissionRepository) Update(
	ctx context.Context,
	tx *gorm.DB,
	permission *entity.Permission,
) error {
	permission.UpdatedAt = time.Now().Unix()

	if tx != nil {
		return tx.WithContext(ctx).Save(&permission).Error
	}

	return r.db.WithContext(ctx).Save(&permission).Error
}

func (r *permissionRepository) Delete(
	ctx context.Context,
	tx *gorm.DB,
	permission *entity.Permission,
) error {
	if tx != nil {
		return tx.WithContext(ctx).Delete(&permission).Error
	}

	return r.db.WithContext(ctx).Delete(&permission).Error
}

func (r *permissionRepository) BulkCreate(
	ctx context.Context,
	tx *gorm.DB,
	permissions []entity.Permission,
) error {
	if tx != nil {
		return tx.WithContext(ctx).Create(&permissions).Error
	}

	return r.db.WithContext(ctx).Create(&permissions).Error
}

func (r *permissionRepository) FindOneByFilter(
	ctx context.Context,
	filter *repository.FindPermissionByFilter,
) (*entity.Permission, error) {
	var data *entity.Permission
	query := r.buildFilter(ctx, nil, filter)

	err := query.First(&data).Error
	return data, err
}

func (r *permissionRepository) FindByFilter(
	ctx context.Context,
	filter *repository.FindPermissionByFilter,
) ([]entity.Permission, error) {
	var data []entity.Permission
	query := r.buildFilter(ctx, nil, filter)

	err := query.Find(&data).Error
	return data, err
}

// -------------------------------------------------------------------------------
func (r *permissionRepository) buildFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindPermissionByFilter,
) *gorm.DB {
	query := r.db.WithContext(ctx)
	if tx != nil {
		query = tx.WithContext(ctx)
	}

	if filter.UserId != nil && *filter.UserId != uuid.Nil {
		query = query.Scopes(findByString[uuid.UUID](*filter.UserId, "user_id"))
	}
	if filter.WorkspaceId != nil && *filter.WorkspaceId != uuid.Nil {
		query = query.Scopes(findByString[uuid.UUID](*filter.WorkspaceId, "workspace_id"))
	}
	if filter.UserWorkspaceId != nil && *filter.UserWorkspaceId != uuid.Nil {
		query = query.Scopes(findByString[uuid.UUID](*filter.UserWorkspaceId, "user_workspace_id"))
	}
	if filter.Permission != nil && *filter.Permission != "" {
		searchText := "%/" + *filter.Permission + "/%"
		query = query.Where("permission LIKE ?", searchText)
	}

	return query
}
