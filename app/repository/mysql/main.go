package mysql_repository

import (
	"gin-boilerplate/app/repository"

	"gorm.io/gorm"
)

type MysqlRepositoryCollections struct {
	MysqlUserRepo repository.UserRepository
}

func RegisterMysqlRepositories(db *gorm.DB) MysqlRepositoryCollections {
	mysqlUserRepo := NewMysqlUserRepository(db)

	return MysqlRepositoryCollections{
		MysqlUserRepo: mysqlUserRepo,
	}
}
