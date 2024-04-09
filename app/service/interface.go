package service

import (
	"context"
	"gin-boilerplate/app/model"
)

type UserService interface {
	Login(ctx context.Context, data *model.LoginRequest) (*model.LoginResponse, error)
	Register(ctx context.Context, data *model.RegisterRequest) (*model.RegisterResponse, error)
	Logout(ctx context.Context, data *model.LogoutRequest) (*model.LogoutResponse, error)
}
