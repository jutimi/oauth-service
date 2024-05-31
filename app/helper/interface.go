package helper

import "gin-boilerplate/app/entity"

type OauthHelper interface {
	GenerateAccessToken(user entity.User) (string, error)
	GenerateRefreshToken(user entity.User) (string, error)
}
