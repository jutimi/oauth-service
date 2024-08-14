package helper

import (
	"context"
	"oauth-server/app/entity"
	"oauth-server/app/repository"
	postgres_repository "oauth-server/app/repository/postgres"
	"oauth-server/config"
	"oauth-server/package/database"
	"oauth-server/package/errors"
	"oauth-server/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jutimi/grpc-service/workspace"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type oauthHelper struct {
	postgresRepo postgres_repository.PostgresRepositoryCollections
}

func NewOauthHelper(
	postgresRepo postgres_repository.PostgresRepositoryCollections,
) OauthHelper {
	return &oauthHelper{
		postgresRepo: postgresRepo,
	}
}

func (h *oauthHelper) GenerateUserToken(
	ctx context.Context,
	user *entity.User,
	tokenType string,
) (string, error) {
	var key string
	conf := config.GetConfiguration().Jwt

	claims := &utils.UserPayload{
		Id:    user.Id,
		Scope: utils.USER_SCOPE,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: conf.Issuer,
		},
	}

	switch tokenType {
	case utils.ACCESS_TOKEN:
		key = conf.UserAccessTokenKey
		claims.RegisteredClaims.ExpiresAt = utils.GenerateExpireTime(utils.USER_ACCESS_TOKEN_IAT)
	case utils.REFRESH_TOKEN:
		key = conf.UserRefreshTokenKey
		claims.RegisteredClaims.ExpiresAt = utils.GenerateExpireTime(utils.USER_REFRESH_TOKEN_IAT)
	}

	token, err := utils.GenerateToken(claims, key)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (h *oauthHelper) GenerateWorkspaceToken(
	ctx context.Context,
	userWorkspace *workspace.UserWorkspaceDetail,
	tokenType string,
) (string, error) {
	var key string
	conf := config.GetConfiguration().Jwt

	userId, err := utils.ConvertStringToUUID(userWorkspace.UserId)
	if err != nil {
		return "", err
	}
	workspaceId, err := utils.ConvertStringToUUID(userWorkspace.WorkspaceId)
	if err != nil {
		return "", err
	}
	userWorkspaceId, err := utils.ConvertStringToUUID(userWorkspace.Id)
	if err != nil {
		return "", err
	}

	claims := &utils.WorkspacePayload{
		Id:              userId,
		Scope:           utils.WORKSPACE_SCOPE,
		WorkspaceId:     workspaceId,
		UserWorkspaceId: userWorkspaceId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: conf.Issuer,
		},
	}
	switch tokenType {
	case utils.ACCESS_TOKEN:
		claims.RegisteredClaims.ExpiresAt = utils.GenerateExpireTime(utils.Workspace_ACCESS_TOKEN_IAT)
		key = conf.WorkspaceAccessTokenKey
	case utils.REFRESH_TOKEN:
		claims.RegisteredClaims.ExpiresAt = utils.GenerateExpireTime(utils.Workspace_REFRESH_TOKEN_IAT)
		key = conf.WorkspaceRefreshTokenKey
	}

	token, err := utils.GenerateToken(claims, key)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (h *oauthHelper) ValidateRefreshToken(ctx context.Context, data *ValidateRefreshTokenParams) error {
	// Check user oauth
	userOAuth, err := h.postgresRepo.OAuthRepo.FindOneByFilter(ctx, &repository.FindOAuthByFilter{
		Token: &data.RefreshToken,
		Scope: &data.Scope,
	})

	if err != nil || userOAuth.Status != entity.OAuthStatusActive {
		return errors.New(errors.ErrCodeUserNotFound)
	}
	if userOAuth.UserId != data.UserId {
		return errors.New(errors.ErrCodeUnauthorized)
	}
	currentTime := time.Now().Unix()
	if currentTime > userOAuth.ExpireAt {
		return errors.New(errors.ErrCodeRefreshTokenExpired)
	}

	return nil
}

func (h *oauthHelper) DeActiveToken(ctx context.Context, data *DeActiveTokenParams) error {
	// Deactivate User OAuth
	tx := database.BeginPostgresTransaction()
	if tx.Error != nil {
		return errors.New(errors.ErrCodeInternalServerError)
	}

	// Find User OAuth
	filter := &repository.FindOAuthByFilter{
		Token: &data.Token,
		Scope: &data.Scope,
	}
	userOAuth, err := h.postgresRepo.OAuthRepo.FindOneByFilterForUpdate(ctx, &repository.FindByFilterForUpdateParams{
		Filter:     filter,
		LockOption: clause.LockingOptionsNoWait,
		Tx:         tx,
	})
	if err != nil {
		return errors.New(errors.ErrCodeInternalServerError)
	}

	userOAuth.Status = entity.OAuthStatusInactive
	if err := h.postgresRepo.OAuthRepo.Update(ctx, tx, userOAuth); err != nil {
		tx.Rollback()
		return errors.New(errors.ErrCodeInternalServerError)
	}

	if err := tx.Commit().Error; err != nil {
		return errors.New(errors.ErrCodeInternalServerError)
	}

	return nil
}

func (h *oauthHelper) ActiveToken(ctx context.Context, data *ActiveTokenParams) error {
	tx := database.BeginPostgresTransaction()
	if tx.Error != nil {
		return errors.New(errors.ErrCodeInternalServerError)
	}

	filter := &repository.FindOAuthByFilter{
		UserId: &data.UserId,
		Scope:  &data.Scope,
	}
	userOAuth, err := h.postgresRepo.OAuthRepo.FindOneByFilterForUpdate(ctx, &repository.FindByFilterForUpdateParams{
		Filter:     filter,
		LockOption: clause.LockingOptionsNoWait,
		Tx:         tx,
	})
	if err == gorm.ErrRecordNotFound {
		userOAuth = entity.NewOAuth()
		userOAuth.UserId = data.UserId
		userOAuth.Status = entity.OAuthStatusActive
	}

	userOAuth.Token = data.Token
	userOAuth.ExpireAt = time.Now().Add(time.Duration(data.TokenIAT) * time.Second).Unix()
	userOAuth.LoginAt = time.Now().Unix()
	if err := h.postgresRepo.OAuthRepo.Update(ctx, tx, userOAuth); err != nil {
		tx.Rollback()
		return errors.New(errors.ErrCodeInternalServerError)
	}

	if err := tx.Commit().Error; err != nil {
		return errors.New(errors.ErrCodeInternalServerError)
	}

	return nil
}
