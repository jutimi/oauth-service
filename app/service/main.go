package service

import "gin-boilerplate/app/repository"

type ServiceCollections struct {
	UserSvc UserService
}

func RegisterServices(repo repository.RepositoryCollections) ServiceCollections {
	databaseSvc := NewDatabaseService(
		repo.MysqlUserRepo,
	)
	userSvc := NewUserService(databaseSvc)

	return ServiceCollections{
		UserSvc: userSvc,
	}
}
