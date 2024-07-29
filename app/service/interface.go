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
	AddUserWSPermission(ctx context.Context, data *model.AddUserWSPermissionRequest) (*model.AddUserWSPermissionResponse, error)
	RevokeUserWSPermission(ctx context.Context, data *model.RevokeUserWSPermissionRequest) (*model.RevokeUserWSPermissionResponse, error)
	GetPermissions(ctx context.Context) (*model.GetPermissionsResponse, error)
}
