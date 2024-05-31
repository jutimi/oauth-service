package service

import (
	"gin-boilerplate/app/helper"
	mysql_repository "gin-boilerplate/app/repository/mysql"
)

type ServiceCollections struct {
	UserSvc  UserService
	OAuthSvc OAuthService
}

func RegisterServices(
	helpers helper.HelperCollections,

	mysqlRepo mysql_repository.MysqlRepositoryCollections,
) ServiceCollections {
	databaseSvc := NewDatabaseService(mysqlRepo)
	userSvc := NewUserService(helpers, databaseSvc)
	oauthSvc := NewOAuthService(helpers, databaseSvc)

	return ServiceCollections{
		UserSvc:  userSvc,
		OAuthSvc: oauthSvc,
	}
}
