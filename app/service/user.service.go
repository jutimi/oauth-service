package service

import (
	"context"
	"fmt"
	"gin-boilerplate/app/model"
	"gin-boilerplate/app/repository"
	"gin-boilerplate/config"
	"gin-boilerplate/package/errors"
	logger "gin-boilerplate/package/log"
	"gin-boilerplate/utils"
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
	conf := config.GetConfiguration().Jwt

	user, err := s.databaseSvc.FindUserByFilter(&repository.FindUserByFilter{
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
		return nil, errors.New(errors.ErrCodeUserNotFound)
	}

	payload := &utils.UserPayload{
		ID: user.ID,
	}
	accessToken, err := utils.GenerateToken(payload, conf.UserAccessTokenKey, 15*60)
	if err != nil {
		logger.Println(logger.LogPrintln{
			FileName:  "app/service/user.service.go",
			FuncName:  "Login",
			TraceData: fmt.Sprintf("%s/%s", data.Email, data.PhoneNumber),
			Msg:       "GenerateAccessToken - " + err.Error(),
		})
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}
	refreshToken, err := utils.GenerateToken(payload, conf.UserAccessTokenKey, 24*60*60)
	if err != nil {
		logger.Println(logger.LogPrintln{
			FileName:  "app/service/user.service.go",
			FuncName:  "Login",
			TraceData: fmt.Sprintf("%s/%s", data.Email, data.PhoneNumber),
			Msg:       "GenerateRefreshToken - " + err.Error(),
		})
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return &model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *userService) Register(ctx context.Context, data *model.RegisterRequest) (*model.RegisterResponse, error) {
	existedUser, err := s.databaseSvc.FindUsersByFilter(&repository.FindUserByFilter{
		PhoneNumber: data.PhoneNumber,
		Email:       data.Email,
	})
	if err != nil {
		logger.Println(logger.LogPrintln{
			FileName:  "app/service/user.service.go",
			FuncName:  "Register",
			TraceData: fmt.Sprintf("%s/%s", data.Email, data.PhoneNumber),
			Msg:       "databaseSvc FindUsersByFilter - " + err.Error(),
		})
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}
	if len(existedUser) > 0 {
		return nil, errors.New(errors.ErrCodeUserExisted)
	}

	return &model.RegisterResponse{}, nil
}
func (s *userService) Logout(ctx context.Context, data *model.LogoutRequest) (*model.LogoutResponse, error) {
	return nil, nil
}
