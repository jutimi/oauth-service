package service

import (
	"gin-boilerplate/app/entity"
	"gin-boilerplate/app/repository"

	"gorm.io/gorm"
)

type databaseService struct {
	mysqlUseRepo repository.UserRepository
}

func NewDatabaseService(
	mysqlUseRepo repository.UserRepository,
) DatabaseService {
	return &databaseService{
		mysqlUseRepo,
	}
}

func (s *databaseService) CreateUser(user *entity.User) error {
	return s.mysqlUseRepo.CreateUser(user)
}
func (s *databaseService) UpdateUser(user *entity.User) error {
	return s.mysqlUseRepo.UpdateUser(user)
}
func (s *databaseService) DeleteUser(user *entity.User) error {
	return s.mysqlUseRepo.DeleteUser(user)
}
func (s *databaseService) NewUser() *entity.User {
	return s.mysqlUseRepo.NewUser()
}
func (s *databaseService) BulkCreateUser(users []entity.User) error {
	return s.mysqlUseRepo.BulkCreateUser(users)
}
func (s *databaseService) FindUserByFilter(filter *repository.FindUserByFilter) (*entity.User, error) {
	return s.mysqlUseRepo.FindUserByFilter(filter)
}

func (s *databaseService) FindUsersByFilter(filter *repository.FindUserByFilter) ([]entity.User, error) {
	return s.mysqlUseRepo.FindUsersByFilter(filter)
}

func (s *databaseService) NewUserTransaction() *gorm.DB {
	return s.mysqlUseRepo.NewUserTransaction()
}
