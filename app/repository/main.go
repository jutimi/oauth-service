package repository

import (
	"gorm.io/gorm"
)

type RepositoryCollections struct {
	MysqlUserRepo UserRepository
}

func RegisterMysqlRepositories(db *gorm.DB) RepositoryCollections {
	mysqlUserRepo := NewMysqlUserRepository(db)

	return RepositoryCollections{
		MysqlUserRepo: mysqlUserRepo,
	}
}
