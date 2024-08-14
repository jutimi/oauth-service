package service

import (
	"context"
	"oauth-server/app/model"
)

type UserService interface {
	Register(ctx context.Context, data *model.RegisterRequest) (*model.RegisterResponse, error)
}

type OAuthService interface {
	RefreshToken(ctx context.Context, data *model.RefreshTokenRequest) (*model.RefreshTokenResponse, error)
	Login(ctx context.Context, data interface{}) (interface{}, error)
	Logout(ctx context.Context, data interface{}) (interface{}, error)
}

type PermissionService interface {
	AddUserWorkspacePermission(ctx context.Context, data *model.AddUserWorkspacePermissionRequest) (*model.AddUserWorkspacePermissionResponse, error)
	RevokeUserWorkspacePermission(ctx context.Context, data *model.RevokeUserWorkspacePermissionRequest) (*model.RevokeUserWorkspacePermissionResponse, error)
	GetPermissions(ctx context.Context) (*model.GetPermissionsResponse, error)
}
