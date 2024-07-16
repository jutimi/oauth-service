package service

import (
	"context"
	"oauth-server/app/entity"
	"oauth-server/app/helper"
	"oauth-server/app/model"
	"oauth-server/app/repository"
	postgres_repository "oauth-server/app/repository/postgres"
	"oauth-server/config"
	"oauth-server/package/errors"
	"oauth-server/utils"
	"time"
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
