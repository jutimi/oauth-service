package helper

import (
	"context"
	"fmt"
	"oauth-server/app/entity"
	"oauth-server/app/repository"
	postgres_repository "oauth-server/app/repository/postgres"
	"oauth-server/config"
	"oauth-server/package/database"
	"oauth-server/package/errors"
	logger "oauth-server/package/log"
	"oauth-server/utils"
	"time"

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

func (h *oauthHelper) GenerateUserToken(user *entity.User, tokenType string) (string, error) {
	var token string
	var err error
	conf := config.GetConfiguration().Jwt

	payload := &utils.UserPayload{
		ID:    user.ID,
		Scope: utils.USER_SCOPE,
	}

	switch tokenType {
	case utils.ACCESS_TOKEN:
		token, err = utils.GenerateToken(payload, conf.UserAccessTokenKey, utils.USER_ACCESS_TOKEN_IAT)
		if err != nil {
			logger.Println(logger.LogPrintln{
				FileName:  "app/helper/oauth.helper.go",
				FuncName:  "GenerateAccessToken",
				TraceData: fmt.Sprintf("%s/%s", *user.Email, *user.PhoneNumber),
				Msg:       err.Error(),
			})
			return token, err
		}
	case utils.REFRESH_TOKEN:
		token, err = utils.GenerateToken(payload, conf.UserRefreshTokenKey, utils.USER_REFRESH_TOKEN_IAT)
		if err != nil {
			logger.Println(logger.LogPrintln{
				FileName:  "app/service/user.service.go",
				FuncName:  "GenerateRefreshToken",
				TraceData: fmt.Sprintf("%s/%s", *user.Email, *user.PhoneNumber),
				Msg:       err.Error(),
			})
			return token, err
		}
	}

	return token, nil
}

func (h *oauthHelper) GenerateWSToken(userWS *workspace.UserWorkspaceDetail, tokenType string) (string, error) {
	var token string
	var err error
	conf := config.GetConfiguration().Jwt

	userId, err := utils.ConvertStringToUUID(userWS.UserId)
	if err != nil {
		return "", err
	}
	workspaceId, err := utils.ConvertStringToUUID(userWS.WorkspaceId)
	if err != nil {
		return "", err
	}
	userWorkspaceId, err := utils.ConvertStringToUUID(userWS.Id)
	if err != nil {
		return "", err
	}

	payload := &utils.WorkspacePayload{
		ID:              userId,
		Scope:           utils.WORKSPACE_SCOPE,
		WorkspaceID:     workspaceId,
		UserWorkspaceID: userWorkspaceId,
	}

	switch tokenType {
	case utils.ACCESS_TOKEN:
		token, err = utils.GenerateToken(payload, conf.WSAccessTokenKey, utils.WS_ACCESS_TOKEN_IAT)
		if err != nil {
			return token, err
		}
	case utils.REFRESH_TOKEN:
		token, err = utils.GenerateToken(payload, conf.WSRefreshTokenKey, utils.WS_REFRESH_TOKEN_IAT)
		if err != nil {
			return token, err
		}
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
	if userOAuth.UserID != data.UserID {
		return errors.New(errors.ErrCodeUnauthorized)
	}
	currentTime := time.Now().Unix()
	if currentTime > userOAuth.ExpireAt {
		return errors.New(errors.ErrCodeTokenExpired)
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
		UserID: &data.UserID,
		Scope:  &data.Scope,
	}
	userOAuth, err := h.postgresRepo.OAuthRepo.FindOneByFilterForUpdate(ctx, &repository.FindByFilterForUpdateParams{
		Filter:     filter,
		LockOption: clause.LockingOptionsNoWait,
		Tx:         tx,
	})
	if err == gorm.ErrRecordNotFound {
		userOAuth = entity.NewOAuth()
		userOAuth.UserID = data.UserID
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
