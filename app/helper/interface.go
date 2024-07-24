package helper

import (
	"context"
	"oauth-server/app/entity"

	"github.com/jutimi/grpc-service/workspace"
)

type OauthHelper interface {
	GenerateUserToken(user *entity.User, tokenType string) (string, error)
	GenerateWSToken(userWS *workspace.UserWorkspaceDetail, tokenType string) (string, error)
	ValidateRefreshToken(ctx context.Context, data *ValidateRefreshTokenParams) error
	DeActiveToken(ctx context.Context, data *DeActiveTokenParams) error
	ActiveToken(ctx context.Context, data *ActiveTokenParams) error
}

type UserHelper interface {
	CreateUser(ctx context.Context, user *entity.User) error
}
