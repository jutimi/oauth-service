package helper

import (
	"context"
	"oauth-server/app/entity"
)

type OauthHelper interface {
	GenerateAccessToken(user entity.User) (string, error)
	GenerateRefreshToken(user entity.User) (string, error)
}

type UserHelper interface {
	CreateUser(ctx context.Context, user *entity.User) error
}
