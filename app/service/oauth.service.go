package service

import (
	"context"
	"oauth-server/app/entity"
	"oauth-server/app/helper"
	"oauth-server/app/model"
	"oauth-server/app/repository"
	postgres_repository "oauth-server/app/repository/postgres"
	"oauth-server/config"
	client_grpc "oauth-server/grpc/client"
	"oauth-server/package/database"
	"oauth-server/package/errors"
	"oauth-server/utils"
	"time"

	"gorm.io/gorm"
)

type oAuthService struct {
	helpers      helper.HelperCollections
	postgresRepo postgres_repository.PostgresRepositoryCollections
	clientGRPC   client_grpc.ClientGRPCCollection
}

func NewOAuthService(
	helpers helper.HelperCollections,
	postgresRepo postgres_repository.PostgresRepositoryCollections,
	clientGRPC client_grpc.ClientGRPCCollection,
) OAuthService {
	return &oAuthService{
		helpers:      helpers,
		postgresRepo: postgresRepo,
		clientGRPC:   clientGRPC,
	}
}

func (s *oAuthService) RefreshToken(ctx context.Context, data *model.RefreshTokenRequest) (*model.RefreshTokenResponse, error) {
	conf := config.GetConfiguration().Jwt

	// Verify refresh token
	claims, err := utils.VerifyToken(data.RefreshToken, conf.UserRefreshTokenKey)
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}
	userPayload, ok := claims.(*utils.UserPayload)
	if !ok {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	// Check user oauth
	userOAuth, err := s.postgresRepo.OAuthRepo.FindOneByFilter(ctx, nil, &repository.FindOAuthByFilter{
		Token: &data.RefreshToken,
	})

	if err != nil || userOAuth.Status != entity.OAuthStatusActive {
		return nil, errors.New(errors.ErrCodeUserNotFound)
	}
	if userOAuth.UserID != userPayload.ID {
		return nil, errors.New(errors.ErrCodeUnauthorized)
	}
	currentTime := time.Now().Unix()
	if currentTime > userOAuth.ExpireAt {
		return nil, errors.New(errors.ErrCodeTokenExpired)
	}

	// Check user exit
	user, err := s.postgresRepo.UserRepo.FindOneByFilter(ctx, nil, &repository.FindUserByFilter{
		ID: &userOAuth.UserID,
	})
	if err != nil {
		return nil, errors.New(errors.ErrCodeUserNotFound)
	}

	// Generate new token
	accessToken, err := s.helpers.OauthHelper.GenerateAccessToken(*user)
	if err != nil {
		return nil, err
	}

	return &model.RefreshTokenResponse{
		AccessToken: accessToken,
	}, nil
}

func (s *oAuthService) Login(ctx context.Context, data interface{}) (interface{}, error) {
	var userOAuth *entity.Oauth
	scope, err := utils.GetScopeContext(ctx)
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
		user, err := s.postgresRepo.UserRepo.FindOneByFilter(ctx, nil, &repository.FindUserByFilter{
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
		accessToken, err := s.helpers.OauthHelper.GenerateAccessToken(*user)
		if err != nil || accessToken == "" {
			return nil, errors.New(errors.ErrCodeInternalServerError)
		}
		refreshToken, err := s.helpers.OauthHelper.GenerateRefreshToken(*user)
		if err != nil || refreshToken == "" {
			return nil, errors.New(errors.ErrCodeInternalServerError)
		}

		// Create User OAuth
		tx := database.BeginPostgresTransaction()
		userOAuth, err = s.postgresRepo.OAuthRepo.FindOneByFilter(ctx, tx, &repository.FindOAuthByFilter{
			UserID: &user.ID,
		})
		if err == gorm.ErrRecordNotFound {
			userOAuth = entity.NewOAuth()
			userOAuth.UserID = user.ID
			userOAuth.Status = entity.OAuthStatusActive
		}

		userOAuth.Token = refreshToken
		userOAuth.ExpireAt = time.Now().Add(utils.USER_REFRESH_TOKEN_IAT * time.Second).Unix()
		userOAuth.LoginAt = time.Now().Unix()
		if err := s.postgresRepo.OAuthRepo.Update(ctx, tx, userOAuth); err != nil {
			tx.WithContext(ctx).Rollback()
			return nil, errors.New(errors.ErrCodeInternalServerError)
		}
		tx.WithContext(ctx).Commit()

		return &model.UserLoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}, nil
	case utils.WORKSPACE_SCOPE:
		_, ok := data.(*model.WorkspaceLoginRequest)
		if !ok {
			return nil, errors.New(errors.ErrCodeInternalServerError)
		}

		// Check user workspace exist
	}

	return nil, errors.New(errors.ErrCodeMethodNotSupported)
}

func (s *oAuthService) Logout(ctx context.Context, data interface{}) (interface{}, error) {
	scope, err := utils.GetScopeContext(ctx)
	if err != nil {
		return nil, err
	}

	token := ""
	switch scope {
	case utils.USER_SCOPE:
		form, ok := data.(*model.UserLogoutRequest)
		if !ok {
			return nil, errors.New(errors.ErrCodeInternalServerError)
		}

		token = form.RefreshToken
	case utils.WORKSPACE_SCOPE:
		form, ok := data.(*model.WorkspaceLogoutRequest)
		if !ok {
			return nil, errors.New(errors.ErrCodeInternalServerError)
		}

		token = form.RefreshToken
	default:
		return nil, errors.New(errors.ErrCodeMethodNotSupported)
	}

	// Find User OAuth
	userOAuth, err := s.postgresRepo.OAuthRepo.FindOneByFilter(ctx, nil, &repository.FindOAuthByFilter{
		Token: &token,
		Scope: &scope,
	})
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	// Deactivate User OAuth
	tx := database.BeginPostgresTransaction()
	userOAuth.Status = entity.OAuthStatusInactive
	if err := s.postgresRepo.OAuthRepo.Update(ctx, tx, userOAuth); err != nil {
		tx.WithContext(ctx).Rollback()
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}
	tx.WithContext(ctx).Commit()

	switch scope {
	case utils.USER_SCOPE:
		return &model.UserLogoutResponse{}, nil
	case utils.WORKSPACE_SCOPE:
		return &model.WorkspaceLogoutResponse{}, nil
	}

	return nil, errors.New(errors.ErrCodeMethodNotSupported)
}
