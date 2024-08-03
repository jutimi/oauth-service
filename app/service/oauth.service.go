package service

import (
	"context"
	"oauth-server/app/helper"
	"oauth-server/app/model"
	"oauth-server/app/repository"
	postgres_repository "oauth-server/app/repository/postgres"
	"oauth-server/config"
	"oauth-server/external/client"
	"oauth-server/package/errors"
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
	scope, err := utils.GetScopeContext[string](ctx, string(utils.SCOPE_CONTEXT_KEY))
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	switch scope {
	case utils.USER_SCOPE:
		// Verify refresh token
		claims, err := utils.VerifyToken(data.RefreshToken, conf.UserRefreshTokenKey)
		if err != nil {
			return nil, errors.New(errors.ErrCodeInternalServerError)
		}
		userPayload, ok := claims.(*utils.UserPayload)
		if !ok {
			return nil, errors.New(errors.ErrCodeInternalServerError)
		}

		if err := s.helpers.OauthHelper.ValidateRefreshToken(ctx, &helper.ValidateRefreshTokenParams{
			RefreshToken: data.RefreshToken,
			Scope:        scope,
			UserID:       userPayload.ID,
		}); err != nil {
			return nil, err
		}

		// Check user exit
		user, err := s.postgresRepo.UserRepo.FindOneByFilter(ctx, &repository.FindUserByFilter{
			ID: &userPayload.ID,
		})
		if err != nil {
			return nil, errors.New(errors.ErrCodeUserNotFound)
		}

		// Generate new token
		accessToken, err := s.helpers.OauthHelper.GenerateUserToken(user, utils.ACCESS_TOKEN)
		if err != nil {
			return nil, err
		}

		return &model.RefreshTokenResponse{
			AccessToken: accessToken,
		}, nil
	case utils.WORKSPACE_SCOPE:
		clientGRPC := client.NewWsClient()
		defer clientGRPC.CloseConn()

		// Verify refresh token
		claims, err := utils.VerifyToken(data.RefreshToken, conf.WSRefreshTokenKey)
		if err != nil {
			return nil, errors.New(errors.ErrCodeInternalServerError)
		}
		wsPayload, ok := claims.(*utils.WorkspacePayload)
		if !ok {
			return nil, errors.New(errors.ErrCodeInternalServerError)
		}

		if err := s.helpers.OauthHelper.ValidateRefreshToken(ctx, &helper.ValidateRefreshTokenParams{
			RefreshToken: data.RefreshToken,
			Scope:        scope,
			UserID:       wsPayload.ID,
		}); err != nil {
			return nil, err
		}

		// Check user workspace exit
		id := wsPayload.UserWorkspaceID.String()
		isActive := true
		userWS, err := clientGRPC.GetUserWSByFilter(ctx, &workspace.GetUserWorkspaceByFilterParams{
			Id:       &id,
			IsActive: &isActive,
		})
		if err != nil {
			return nil, err
		}

		// Generate new token
		accessToken, err := s.helpers.OauthHelper.GenerateWSToken(userWS.Data, utils.ACCESS_TOKEN)
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
	scope, err := utils.GetScopeContext[string](ctx, string(utils.SCOPE_CONTEXT_KEY))
	if err != nil {
		return nil, err
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
		accessToken, err := s.helpers.OauthHelper.GenerateUserToken(user, utils.ACCESS_TOKEN)
		if err != nil || accessToken == "" {
			return nil, errors.New(errors.ErrCodeInternalServerError)
		}
		refreshToken, err := s.helpers.OauthHelper.GenerateUserToken(user, utils.REFRESH_TOKEN)
		if err != nil || refreshToken == "" {
			return nil, errors.New(errors.ErrCodeInternalServerError)
		}

		if err := s.helpers.OauthHelper.ActiveToken(ctx, &helper.ActiveTokenParams{
			Token:    refreshToken,
			Scope:    scope,
			UserID:   user.ID,
			TokenIAT: utils.USER_REFRESH_TOKEN_IAT,
		}); err != nil {
			return nil, err
		}

		return &model.UserLoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}, nil
	case utils.WORKSPACE_SCOPE:
		clientGRPC := client.NewWsClient()
		defer clientGRPC.CloseConn()

		form, ok := data.(*model.WorkspaceLoginRequest)
		if !ok {
			return nil, errors.New(errors.ErrCodeInternalServerError)
		}

		// Check user workspace exist
		userPayload, err := utils.GetScopeContext[*utils.UserPayload](ctx, string(utils.USER_CONTEXT_KEY))
		if err != nil {
			return nil, errors.New(errors.ErrCodeInternalServerError)
		}
		userId := userPayload.ID.String()
		userWS, err := clientGRPC.GetUserWSByFilter(ctx, &workspace.GetUserWorkspaceByFilterParams{
			WorkspaceId: &form.WorkspaceID,
			UserId:      &userId,
		})
		if err != nil {
			return nil, err
		}

		// Generate token
		accessToken, err := s.helpers.OauthHelper.GenerateWSToken(userWS.GetData(), utils.ACCESS_TOKEN)
		if err != nil || accessToken == "" {
			return nil, errors.New(errors.ErrCodeInternalServerError)
		}
		refreshToken, err := s.helpers.OauthHelper.GenerateWSToken(userWS.GetData(), utils.REFRESH_TOKEN)
		if err != nil || refreshToken == "" {
			return nil, errors.New(errors.ErrCodeInternalServerError)
		}

		// Create User OAuth
		if err := s.helpers.OauthHelper.ActiveToken(ctx, &helper.ActiveTokenParams{
			Token:    refreshToken,
			Scope:    scope,
			UserID:   userPayload.ID,
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
	scope, err := utils.GetScopeContext[string](ctx, string(utils.SCOPE_CONTEXT_KEY))
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
