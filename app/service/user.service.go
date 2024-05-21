package service

import (
	"context"
	"fmt"
	"gin-boilerplate/app/model"
	"gin-boilerplate/app/repository"
	"gin-boilerplate/package/errors"
	logger "gin-boilerplate/package/log"
)

type userService struct {
	databaseSvc DatabaseService
}

func NewUserService(
	databaseSvc DatabaseService,
) UserService {
	return &userService{
		databaseSvc: databaseSvc,
	}
}

func (s *userService) Login(ctx context.Context, data *model.LoginRequest) (*model.LoginResponse, error) {
	existedUser, err := s.databaseSvc.FindUsersByFilter(&repository.FindUserByFilter{
		PhoneNumber: data.PhoneNumber,
		Email:       data.Email,
	})
	if err != nil {
		logger.Println(logger.LogPrintln{
			FileName:  "app/service/user.service.go",
			FuncName:  "Login",
			TraceData: fmt.Sprintf("%s/%s", data.Email, data.PhoneNumber),
			Msg:       "databaseSvc FindUserByFilter - " + err.Error(),
		})
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	if len(existedUser) > 0 {
		return nil, errors.New(errors.ErrCodeUserExisted)
	}

	return nil, nil
}
func (s *userService) Register(ctx context.Context, data *model.RegisterRequest) (*model.RegisterResponse, error) {
	return nil, nil
}
func (s *userService) Logout(ctx context.Context, data *model.LogoutRequest) (*model.LogoutResponse, error) {
	return nil, nil
}
