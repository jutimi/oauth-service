package postgres_repository

import (
	"oauth-server/app/repository"

	"gorm.io/gorm"
)

type PostgresRepositoryCollections struct {
	PostgresUserRepo  repository.UserRepository
	PostgresOAuthRepo repository.OAuthRepository
}

func RegisterPostgresRepositories(db *gorm.DB) PostgresRepositoryCollections {
	postgresUserRepo := NewPostgresUserRepository(db)
	postgresOAuthRepo := NewPostgresOAuthRepository(db)

	return PostgresRepositoryCollections{
		PostgresUserRepo:  postgresUserRepo,
		PostgresOAuthRepo: postgresOAuthRepo,
	}
}
