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
