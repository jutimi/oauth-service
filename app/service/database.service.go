package service

import (
	"context"
	"gin-boilerplate/app/entity"
	"gin-boilerplate/app/repository"
	mysql_repository "gin-boilerplate/app/repository/mysql"

	"gorm.io/gorm"
)

type databaseService struct {
	mysqlRepo mysql_repository.MysqlRepositoryCollections
}

func NewDatabaseService(
	mysqlRepo mysql_repository.MysqlRepositoryCollections,
) DatabaseService {
	return &databaseService{
		mysqlRepo: mysqlRepo,
	}
}

func (s *databaseService) CreateUser(
	ctx context.Context,
	tx *gorm.DB,
	user *entity.User,
) error {
	return s.mysqlRepo.MysqlUserRepo.CreateUser(ctx, tx, user)
}
func (s *databaseService) UpdateUser(
	ctx context.Context,
	tx *gorm.DB,
	user *entity.User,
) error {
	return s.mysqlRepo.MysqlUserRepo.UpdateUser(ctx, tx, user)
}
func (s *databaseService) DeleteUser(
	ctx context.Context,
	tx *gorm.DB,
	user *entity.User,
) error {
	return s.mysqlRepo.MysqlUserRepo.DeleteUser(ctx, tx, user)
}

func (s *databaseService) BulkCreateUser(
	ctx context.Context,
	tx *gorm.DB,
	users []entity.User,
) error {
	return s.mysqlRepo.MysqlUserRepo.BulkCreateUser(ctx, tx, users)
}
func (s *databaseService) FindUserByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindUserByFilter,
) (*entity.User, error) {
	return s.mysqlRepo.MysqlUserRepo.FindUserByFilter(ctx, tx, filter)
}

func (s *databaseService) FindUsersByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindUserByFilter,
) ([]entity.User, error) {
	return s.mysqlRepo.MysqlUserRepo.FindUsersByFilter(ctx, tx, filter)
}

func (s *databaseService) CreateOAuth(
	ctx context.Context,
	tx *gorm.DB,
	oauth *entity.Oauth,
) error {
	return s.mysqlRepo.MysqlOAuthRepo.CreateOAuth(ctx, tx, oauth)
}

func (s *databaseService) UpdateOAuth(
	ctx context.Context,
	tx *gorm.DB,
	oauth *entity.Oauth,
) error {
	return s.mysqlRepo.MysqlOAuthRepo.UpdateOAuth(ctx, tx, oauth)
}

func (s *databaseService) FindOAuthByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindOAuthByFilter,
) (*entity.Oauth, error) {
	return s.mysqlRepo.MysqlOAuthRepo.FindOAuthByFilter(ctx, tx, filter)
}
