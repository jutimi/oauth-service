package service

import (
	mysql_repository "gin-boilerplate/app/repository/mysql"
)

type ServiceCollections struct {
	UserSvc UserService
}

func RegisterServices(
	mysqlRepo mysql_repository.MysqlRepositoryCollections,
) ServiceCollections {
	databaseSvc := NewDatabaseService(
		mysqlRepo.MysqlUserRepo,
	)
	userSvc := NewUserService(databaseSvc)

	return ServiceCollections{
		UserSvc: userSvc,
	}
}
