package helper

import (
	postgres_repository "oauth-server/app/repository/postgres"
)

type HelperCollections struct {
	OauthHelper OauthHelper
}

func RegisterHelpers(
	postgresRepo postgres_repository.PostgresRepositoryCollections,
) HelperCollections {
	oauthHelper := NewOauthHelper()

	return HelperCollections{
		OauthHelper: oauthHelper,
	}
}
