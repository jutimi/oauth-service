package helper

import (
	postgres_repository "oauth-server/app/repository/postgres"
)

type HelperCollections struct {
	OauthHelper      OauthHelper
	UserHelper       UserHelper
	PermissionHelper PermissionHelper
}

func RegisterHelpers(
	postgresRepo postgres_repository.PostgresRepositoryCollections,
) HelperCollections {
	return HelperCollections{
		OauthHelper:      NewOauthHelper(postgresRepo),
		UserHelper:       NewUserHelper(postgresRepo),
		PermissionHelper: NewPermissionHelper(postgresRepo),
	}
}
