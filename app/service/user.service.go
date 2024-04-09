package service

import (
	"context"
	"gin-boilerplate/app/model"
)

type userService struct {
}

func NewUserService() UserService {
	return &userService{}
}

func (s *userService) Login(ctx context.Context, data *model.LoginRequest) (*model.LoginResponse, error) {
	return nil, nil
}
func (s *userService) Register(ctx context.Context, data *model.RegisterRequest) (*model.RegisterResponse, error) {
	return nil, nil
}
func (s *userService) Logout(ctx context.Context, data *model.LogoutRequest) (*model.LogoutResponse, error) {
	return nil, nil
}
