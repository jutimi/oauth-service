package mysql_repository

import (
	"gorm.io/gorm"
)

type MysqlRepositoryCollections struct {
}

func RegisterMysqlRepositories(db *gorm.DB) MysqlRepositoryCollections {

	return MysqlRepositoryCollections{}
}
