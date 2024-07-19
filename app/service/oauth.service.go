package service

import (
	"context"
	"oauth-server/app/entity"
	"oauth-server/app/helper"
	"oauth-server/app/model"
	"oauth-server/app/repository"
	postgres_repository "oauth-server/app/repository/postgres"
	"oauth-server/config"
	"oauth-server/package/database"
	"oauth-server/package/errors"
	"oauth-server/utils"
	"time"

	"gorm.io/gorm"
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

func (s *userService) Login(ctx context.Context, data *model.LoginRequest) (*model.LoginResponse, error) {
	var userOAuth *entity.Oauth

	// Check user exit
	user, err := s.postgresRepo.UserRepo.FindOneByFilter(ctx, nil, &repository.FindUserByFilter{
		PhoneNumber: &data.PhoneNumber,
		Email:       &data.Email,
	})
	if err != nil {
		return nil, errors.New(errors.ErrCodeUserNotFound)
	}
	if err := utils.CheckPasswordHash(data.Password, user.Password); err != nil {
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

	return &model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *userService) Logout(ctx context.Context, data *model.LogoutRequest) (*model.LogoutResponse, error) {
	user := ctx.Value(utils.USER_CONTEXT_KEY).(entity.User)

	// Find User OAuth
	userOAuth, err := s.postgresRepo.OAuthRepo.FindOneByFilter(ctx, nil, &repository.FindOAuthByFilter{
		UserID: &user.ID,
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

	return &model.LogoutResponse{}, nil
}
