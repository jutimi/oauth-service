package helper

import "oauth-server/app/entity"

type OauthHelper interface {
	GenerateAccessToken(user entity.User) (string, error)
	GenerateRefreshToken(user entity.User) (string, error)
}
