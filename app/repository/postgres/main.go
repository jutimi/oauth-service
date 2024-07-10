package postgres_repository

import (
	"oauth-server/app/repository"

	"gorm.io/gorm"
)

type PostgresRepositoryCollections struct {
	UserRepo  repository.UserRepository
	OAuthRepo repository.OAuthRepository
}

func RegisterPostgresRepositories(db *gorm.DB) PostgresRepositoryCollections {

	return PostgresRepositoryCollections{
		UserRepo:  NewUserRepository(db),
		OAuthRepo: NewOAuthRepository(db),
	}
}
