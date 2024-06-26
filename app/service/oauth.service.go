package service

import (
	"context"
	"gin-boilerplate/app/entity"
	"gin-boilerplate/app/helper"
	"gin-boilerplate/app/model"
	"gin-boilerplate/app/repository"
	postgres_repository "gin-boilerplate/app/repository/postgres"
	"gin-boilerplate/package/errors"
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

	// Check user oauth
	userOAuth, err := s.postgresRepo.PostgresOAuthRepo.FindOAuthByFilter(ctx, nil, &repository.FindOAuthByFilter{
		Token: data.RefreshToken,
	})
	if err != nil || userOAuth.Status != entity.OAuthStatusActive {
		return nil, errors.New(errors.ErrCodeUserNotFound)
	}
	currentTime := time.Now().Unix()
	if currentTime > userOAuth.ExpireAt {
		return nil, errors.New(errors.ErrCodeTokenExpired)
	}

	// Check user exit
	user, err := s.postgresRepo.PostgresUserRepo.FindUserByFilter(ctx, nil, &repository.FindUserByFilter{
		ID: userOAuth.UserID,
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
