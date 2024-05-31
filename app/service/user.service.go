package service

import (
	"context"
	"fmt"
	"gin-boilerplate/app/entity"
	"gin-boilerplate/app/helper"
	"gin-boilerplate/app/model"
	"gin-boilerplate/app/repository"
	"gin-boilerplate/package/database"
	"gin-boilerplate/package/errors"
	logger "gin-boilerplate/package/log"
	"gin-boilerplate/utils"
	"time"
)

type userService struct {
	helpers     helper.HelperCollections
	databaseSvc DatabaseService
}

func NewUserService(
	helpers helper.HelperCollections,
	databaseSvc DatabaseService,
) UserService {
	return &userService{
		helpers:     helpers,
		databaseSvc: databaseSvc,
	}
}

func (s *userService) Login(ctx context.Context, data *model.LoginRequest) (*model.LoginResponse, error) {
	// Check user exit
	user, err := s.databaseSvc.FindUserByFilter(ctx, nil, &repository.FindUserByFilter{
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

	// Generate token
	accessToken, err := s.helpers.OauthHelper.GenerateAccessToken(*user)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.helpers.OauthHelper.GenerateRefreshToken(*user)
	if err != nil {
		return nil, err
	}

	// Create User OAuth
	tx := database.BeginMysqlTransaction()
	userOAuth := entity.NewOAuth()
	userOAuth.UserID = user.ID
	userOAuth.Token = refreshToken
	userOAuth.Status = entity.OAuthStatusActive
	userOAuth.ExpireAt = time.Now().Add(time.Hour * 24 * 30).Unix()
	if err := s.databaseSvc.CreateOAuth(ctx, tx, userOAuth); err != nil {
		logger.Println(logger.LogPrintln{
			FileName:  "app/service/user.service.go",
			FuncName:  "Login",
			TraceData: fmt.Sprintf("%s/%s", data.Email, data.PhoneNumber),
			Msg:       "GenerateRefreshToken - " + err.Error(),
		})
		tx.Rollback()
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}
	tx.Commit()

	return &model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *userService) Register(ctx context.Context, data *model.RegisterRequest) (*model.RegisterResponse, error) {
	// Check user exited
	existedUser, err := s.databaseSvc.FindUsersByFilter(ctx, nil, &repository.FindUserByFilter{
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

	// Create user
	tx := database.BeginMysqlTransaction()
	user := entity.NewUser()
	user.PhoneNumber = &data.PhoneNumber
	user.Email = &data.Email
	user.Password = data.Password
	if err := s.databaseSvc.CreateUser(ctx, tx, user); err != nil {
		logger.Println(logger.LogPrintln{
			FileName:  "app/service/user.service.go",
			FuncName:  "Register",
			TraceData: fmt.Sprintf("%s/%s", data.Email, data.PhoneNumber),
			Msg:       "databaseSvc CreateUser - " + err.Error(),
		})
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return &model.RegisterResponse{}, nil
}

func (s *userService) Logout(ctx context.Context, data *model.LogoutRequest) (*model.LogoutResponse, error) {
	user := ctx.Value(utils.USER_CONTEXT_KEY).(entity.User)

	// Find User OAuth
	userOAuth, err := s.databaseSvc.FindOAuthByFilter(ctx, nil, &repository.FindOAuthByFilter{
		UserID: user.ID,
	})
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	// Deactivate User OAuth
	tx := database.BeginMysqlTransaction()
	userOAuth.Status = entity.OAuthStatusInactive
	if err := s.databaseSvc.UpdateOAuth(ctx, tx, userOAuth); err != nil {
		tx.Rollback()
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}
	tx.Commit()

	return &model.LogoutResponse{}, nil
}
