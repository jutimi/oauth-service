package helper

import (
	"context"
	"oauth-server/app/entity"

	"github.com/jutimi/grpc-service/workspace"
)

type OauthHelper interface {
	GenerateUserToken(ctx context.Context, user *entity.User, tokenType string) (string, error)
	GenerateWSToken(ctx context.Context, userWS *workspace.UserWorkspaceDetail, tokenType string) (string, error)
	ValidateRefreshToken(ctx context.Context, data *ValidateRefreshTokenParams) error
	DeActiveToken(ctx context.Context, data *DeActiveTokenParams) error
	ActiveToken(ctx context.Context, data *ActiveTokenParams) error
}

type UserHelper interface {
	CreateUser(ctx context.Context, data *CreateUserParams) (*entity.User, error)
	CreateUsers(ctx context.Context, data []CreateUserParams) ([]entity.User, error)
}

type PermissionHelper interface {
	ValidatePermission(ctx context.Context, permission string) error
	GetPermissions(ctx context.Context, permission string) map[string]bool
	GetURLPermission(ctx context.Context, resource, action string) (string, error)
}
