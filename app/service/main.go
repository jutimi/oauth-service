package service

import (
	"gin-boilerplate/app/helper"
	postgres_repository "gin-boilerplate/app/repository/postgres"
)

type ServiceCollections struct {
	UserSvc  UserService
	OAuthSvc OAuthService
}

func RegisterServices(
	helpers helper.HelperCollections,

	postgresRepo postgres_repository.PostgresRepositoryCollections,
) ServiceCollections {
	userSvc := NewUserService(helpers, postgresRepo)
	oauthSvc := NewOAuthService(helpers, postgresRepo)

	return ServiceCollections{
		UserSvc:  userSvc,
		OAuthSvc: oauthSvc,
	}
}
