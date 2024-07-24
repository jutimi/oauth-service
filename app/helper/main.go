package helper

import (
	postgres_repository "oauth-server/app/repository/postgres"
	client_grpc "oauth-server/grpc/client"
)

type HelperCollections struct {
	OauthHelper OauthHelper
}

func RegisterHelpers(
	postgresRepo postgres_repository.PostgresRepositoryCollections,
	clientGRPC client_grpc.ClientGRPCCollection,
) HelperCollections {
	oauthHelper := NewOauthHelper(postgresRepo)

	return HelperCollections{
		OauthHelper: oauthHelper,
	}
}
