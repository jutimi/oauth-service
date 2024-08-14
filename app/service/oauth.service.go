package service

import (
	"context"
	"fmt"
	"oauth-server/app/helper"
	"oauth-server/app/model"
	"oauth-server/app/repository"
	postgres_repository "oauth-server/app/repository/postgres"
	"oauth-server/config"
	"oauth-server/external/client"
	"oauth-server/package/errors"
	logger "oauth-server/package/log"
	"oauth-server/utils"

	"github.com/jutimi/grpc-service/workspace"
)

type oAuthService struct {
	helpers      helper.HelperCollections
	postgresRepo postgres_repository.PostgresRepositoryCollections
}

func NewOAuthService(
	helpers helper.HelperCollections,
	postgresRepo postgres_repository.PostgresRepositoryCollections,
) OAuthService {
	return &oAuthService{
		helpers:      helpers,
		postgresRepo: postgresRepo,
	}
}

func (s *oAuthService) RefreshToken(ctx context.Context, data *model.RefreshTokenRequest) (*model.RefreshTokenResponse, error) {
	conf := config.GetConfiguration().Jwt
	scope, err := utils.GetScopeContext[string](ctx, utils.SCOPE_CONTEXT_KEY)
	if err != nil {
		logger.Println(logger.LogPrintln{
			Ctx:       ctx,
			FileName:  "app/service/oauth.service.go",
			FuncName:  "RefreshToken",
			TraceData: fmt.Sprintf("%+v", data),
			Msg:       fmt.Sprintf("GetScopeContext - %s", err.Error()),
		})
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	switch scope {
	case utils.USER_SCOPE:
		// Verify refresh token
		payload, err := utils.VerifyUserToken(data.RefreshToken, conf.UserRefreshTokenKey)
		if err != nil {
			logger.Println(logger.LogPrintln{
				Ctx:       ctx,
				FileName:  "app/service/oauth.service.go",
				FuncName:  "RefreshToken",
				TraceData: fmt.Sprintf("%+v", data),
				Msg:       fmt.Sprintf("VerifyToken - User - %s", err.Error()),
			})
			return nil, errors.New(errors.ErrCodeInternalServerError)
		}

		if err := s.helpers.OauthHelper.ValidateRefreshToken(ctx, &helper.ValidateRefreshTokenParams{
			RefreshToken: data.RefreshToken,
			Scope:        scope,
			UserId:       payload.Id,
		}); err != nil {
			return nil, err
		}

		// Check user exit
		user, err := s.postgresRepo.UserRepo.FindOneByFilter(ctx, &repository.FindUserByFilter{
			Id: &payload.Id,
		})
		if err != nil {
			return nil, errors.New(errors.ErrCodeUserNotFound)
		}

		// Generate new token
		accessToken, err := s.helpers.OauthHelper.GenerateUserToken(ctx, user, utils.ACCESS_TOKEN)
		if err != nil {
			return nil, err
		}

		return &model.RefreshTokenResponse{
			AccessToken: accessToken,
		}, nil
	case utils.WORKSPACE_SCOPE:
		clientGRPC := client.NewWorkspaceClient()
		defer clientGRPC.CloseConn()

		// Verify refresh token
		payload, err := utils.VerifyWorkspaceToken(data.RefreshToken, conf.WorkspaceRefreshTokenKey)
		if err != nil {
			logger.Println(logger.LogPrintln{
				Ctx:       ctx,
				FileName:  "app/service/oauth.service.go",
				FuncName:  "RefreshToken",
				TraceData: fmt.Sprintf("%+v", data),
				Msg:       fmt.Sprintf("VerifyToken - Workspace - %s", err.Error()),
			})
			return nil, errors.New(errors.ErrCodeInternalServerError)
		}

		if err := s.helpers.OauthHelper.ValidateRefreshToken(ctx, &helper.ValidateRefreshTokenParams{
			RefreshToken: data.RefreshToken,
			Scope:        scope,
			UserId:       payload.Id,
		}); err != nil {
			return nil, err
		}

		// Check user workspace exit
		id := payload.UserWorkspaceId.String()
		isActive := true
		userWorkspace, err := clientGRPC.GetUserWorkspaceByFilter(ctx, &workspace.GetUserWorkspaceByFilterParams{
			Id:       &id,
			IsActive: &isActive,
		})
		if err != nil {
			logger.Println(logger.LogPrintln{
				Ctx:       ctx,
				FileName:  "app/service/oauth.service.go",
				FuncName:  "RefreshToken",
				TraceData: fmt.Sprintf("%+v", data),
				Msg:       fmt.Sprintf("GetUserWorkspaceByFilter - %s", err.Error()),
			})
			return nil, err
		}

		// Generate new token
		accessToken, err := s.helpers.OauthHelper.GenerateWorkspaceToken(ctx, userWorkspace.Data, utils.ACCESS_TOKEN)
		if err != nil {
			return nil, err
		}

		return &model.RefreshTokenResponse{
			AccessToken: accessToken,
		}, nil
	default:
		return nil, errors.New(errors.ErrCodeMethodNotSupported)
	}
}

func (s *oAuthService) Login(ctx context.Context, data interface{}) (interface{}, error) {
	scope, err := utils.GetScopeContext[string](ctx, utils.SCOPE_CONTEXT_KEY)
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	switch scope {
	case utils.USER_SCOPE:
		form, ok := data.(*model.UserLoginRequest)
		if !ok {
			return nil, errors.New(errors.ErrCodeInternalServerError)
		}

		// Check user exit
		user, err := s.postgresRepo.UserRepo.FindOneByFilter(ctx, &repository.FindUserByFilter{
			PhoneNumber: &form.PhoneNumber,
			Email:       &form.Email,
		})
		if err != nil {
			return nil, errors.New(errors.ErrCodeUserNotFound)
		}
		if err := utils.CheckPasswordHash(form.Password, user.Password); err != nil {
			return nil, errors.New(errors.ErrCodeIncorrectPassword)
		}

		// Generate token
		accessToken, err := s.helpers.OauthHelper.GenerateUserToken(ctx, user, utils.ACCESS_TOKEN)
		if err != nil || accessToken == "" {
			return nil, errors.New(errors.ErrCodeInternalServerError)
		}
		refreshToken, err := s.helpers.OauthHelper.GenerateUserToken(ctx, user, utils.REFRESH_TOKEN)
		if err != nil || refreshToken == "" {
			return nil, errors.New(errors.ErrCodeInternalServerError)
		}

		if err := s.helpers.OauthHelper.ActiveToken(ctx, &helper.ActiveTokenParams{
			Token:    refreshToken,
			Scope:    scope,
			UserId:   user.Id,
			TokenIAT: utils.USER_REFRESH_TOKEN_IAT,
		}); err != nil {
			return nil, err
		}

		return &model.UserLoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}, nil
	case utils.WORKSPACE_SCOPE:
		clientGRPC := client.NewWorkspaceClient()
		defer clientGRPC.CloseConn()

		form, ok := data.(*model.WorkspaceLoginRequest)
		if !ok {
			return nil, errors.New(errors.ErrCodeInternalServerError)
		}

		// Check user workspace exist
		userPayload, err := utils.GetGinContext[*utils.UserPayload](ctx, string(utils.USER_CONTEXT_KEY))
		if err != nil {
			return nil, errors.New(errors.ErrCodeInternalServerError)
		}
		userId := userPayload.Id.String()
		userWorkspace, err := clientGRPC.GetUserWorkspaceByFilter(ctx, &workspace.GetUserWorkspaceByFilterParams{
			WorkspaceId: &form.WorkspaceId,
			UserId:      &userId,
		})
		if err != nil {
			return nil, err
		}

		// Generate token
		accessToken, err := s.helpers.OauthHelper.GenerateWorkspaceToken(ctx, userWorkspace.GetData(), utils.ACCESS_TOKEN)
		if err != nil || accessToken == "" {
			return nil, errors.New(errors.ErrCodeInternalServerError)
		}
		refreshToken, err := s.helpers.OauthHelper.GenerateWorkspaceToken(ctx, userWorkspace.GetData(), utils.REFRESH_TOKEN)
		if err != nil || refreshToken == "" {
			return nil, errors.New(errors.ErrCodeInternalServerError)
		}

		// Create User OAuth
		if err := s.helpers.OauthHelper.ActiveToken(ctx, &helper.ActiveTokenParams{
			Token:    refreshToken,
			Scope:    scope,
			UserId:   userPayload.Id,
			TokenIAT: utils.USER_REFRESH_TOKEN_IAT,
		}); err != nil {
			return nil, err
		}

		return &model.WorkspaceLoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}, nil
	}

	return nil, errors.New(errors.ErrCodeMethodNotSupported)
}

func (s *oAuthService) Logout(ctx context.Context, data interface{}) (interface{}, error) {
	scope, err := utils.GetScopeContext[string](ctx, utils.SCOPE_CONTEXT_KEY)
	if err != nil {
		return nil, err
	}

	switch scope {
	case utils.USER_SCOPE:
		form, ok := data.(*model.UserLogoutRequest)
		if !ok {
			return nil, errors.New(errors.ErrCodeInternalServerError)
		}

		if err := s.helpers.OauthHelper.DeActiveToken(ctx, &helper.DeActiveTokenParams{
			Scope: scope,
			Token: form.RefreshToken,
		}); err != nil {
			return nil, err
		}

		return &model.UserLogoutResponse{}, nil
	case utils.WORKSPACE_SCOPE:
		form, ok := data.(*model.WorkspaceLogoutRequest)
		if !ok {
			return nil, errors.New(errors.ErrCodeInternalServerError)
		}

		if err := s.helpers.OauthHelper.DeActiveToken(ctx, &helper.DeActiveTokenParams{
			Scope: scope,
			Token: form.RefreshToken,
		}); err != nil {
			return nil, err
		}

		return &model.WorkspaceLogoutResponse{}, nil
	}

	return nil, errors.New(errors.ErrCodeMethodNotSupported)
}
