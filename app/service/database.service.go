package service

import (
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

func (s *databaseService) CreateUser(user *entity.User) error {
	return s.mysqlRepo.MysqlUserRepo.CreateUser(user)
}
func (s *databaseService) UpdateUser(user *entity.User) error {
	return s.mysqlRepo.MysqlUserRepo.UpdateUser(user)
}
func (s *databaseService) DeleteUser(user *entity.User) error {
	return s.mysqlRepo.MysqlUserRepo.DeleteUser(user)
}
func (s *databaseService) NewUser() *entity.User {
	return s.mysqlRepo.MysqlUserRepo.NewUser()
}
func (s *databaseService) BulkCreateUser(users []entity.User) error {
	return s.mysqlRepo.MysqlUserRepo.BulkCreateUser(users)
}
func (s *databaseService) FindUserByFilter(filter *repository.FindUserByFilter) (*entity.User, error) {
	return s.mysqlRepo.MysqlUserRepo.FindUserByFilter(filter)
}

func (s *databaseService) FindUsersByFilter(filter *repository.FindUserByFilter) ([]entity.User, error) {
	return s.mysqlRepo.MysqlUserRepo.FindUsersByFilter(filter)
}

func (s *databaseService) NewUserTransaction() *gorm.DB {
	return s.mysqlRepo.MysqlUserRepo.NewUserTransaction()
}

func (s *databaseService) NewOAuthTransaction() *gorm.DB {
	return s.mysqlRepo.MysqlOAuthRepo.NewOAuthTransaction()
}

func (s *databaseService) CreateOAuth(oauth *entity.Oauth) error {
	return s.mysqlRepo.MysqlOAuthRepo.CreateOAuth(oauth)
}

func (s *databaseService) UpdateOAuth(oauth *entity.Oauth) error {
	return s.mysqlRepo.MysqlOAuthRepo.UpdateOAuth(oauth)
}

func (s *databaseService) NewOAuth() *entity.Oauth {
	return s.mysqlRepo.MysqlOAuthRepo.NewOAuth()
}

func (s *databaseService) FindOAuthByFilter(filter *repository.FindOAuthByFilter) (*entity.Oauth, error) {
	return s.mysqlRepo.MysqlOAuthRepo.FindOAuthByFilter(filter)
}
