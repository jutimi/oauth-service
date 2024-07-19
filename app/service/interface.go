package service

import (
	"context"
	"oauth-server/app/model"
)

type UserService interface {
	Register(ctx context.Context, data *model.RegisterRequest) (*model.RegisterResponse, error)
}

type OAuthService[
	LoginParams model.UserLoginRequest | model.WorkspaceLoginRequest,
	LoginResponse model.UserLoginResponse | model.WorkspaceLoginResponse,
	LogoutParams model.UserLogoutRequest | model.WorkspaceLogoutRequest,
	LogoutResponse model.UserLogoutResponse | model.WorkspaceLogoutResponse,
] interface {
	RefreshToken(ctx context.Context, data *model.RefreshTokenRequest) (*model.RefreshTokenResponse, error)
	Login(ctx context.Context, data *LoginParams) (*LoginResponse, error)
	Logout(ctx context.Context, data *LogoutParams) (*LoginResponse, error)
}
