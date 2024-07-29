package service

import (
	"oauth-server/app/helper"
	postgres_repository "oauth-server/app/repository/postgres"
)

type ServiceCollections struct {
	UserSvc       UserService
	OAuthSvc      OAuthService
	PermissionSvc PermissionService
}

func RegisterServices(
	helpers helper.HelperCollections,
	postgresRepo postgres_repository.PostgresRepositoryCollections,
) ServiceCollections {
	return ServiceCollections{
		UserSvc:       NewUserService(helpers, postgresRepo),
		OAuthSvc:      NewOAuthService(helpers, postgresRepo),
		PermissionSvc: NewPermissionService(helpers, postgresRepo),
	}
}
