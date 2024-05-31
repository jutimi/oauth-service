package helper

import (
	"fmt"
	"gin-boilerplate/app/entity"
	"gin-boilerplate/config"
	"gin-boilerplate/package/errors"
	logger "gin-boilerplate/package/log"
	"gin-boilerplate/utils"
)

type oauthHelper struct {
}

func NewOauthHelper() OauthHelper {
	return &oauthHelper{}
}

const (
	userAccessTokenIAT  = 15 * 60           // 15 minutes
	userRefreshTokenIAT = 30 * 24 * 60 * 60 // 30 days
)

func (h *oauthHelper) GenerateAccessToken(user entity.User) (string, error) {
	conf := config.GetConfiguration().Jwt

	payload := &utils.UserPayload{
		ID: user.ID,
	}
	accessToken, err := utils.GenerateToken(payload, conf.UserAccessTokenKey, userAccessTokenIAT)
	if err != nil {
		logger.Println(logger.LogPrintln{
			FileName:  "app/helper/oauth.helper.go",
			FuncName:  "GenerateAccessToken",
			TraceData: fmt.Sprintf("%s/%s", *user.Email, *user.PhoneNumber),
			Msg:       err.Error(),
		})
		return "", errors.New(errors.ErrCodeInternalServerError)
	}

	return accessToken, nil
}

func (h *oauthHelper) GenerateRefreshToken(user entity.User) (string, error) {
	conf := config.GetConfiguration().Jwt

	payload := &utils.UserPayload{
		ID: user.ID,
	}
	refreshToken, err := utils.GenerateToken(payload, conf.UserAccessTokenKey, userRefreshTokenIAT)
	if err != nil {
		logger.Println(logger.LogPrintln{
			FileName:  "app/service/user.service.go",
			FuncName:  "GenerateRefreshToken",
			TraceData: fmt.Sprintf("%s/%s", *user.Email, *user.PhoneNumber),
			Msg:       err.Error(),
		})
		return "", errors.New(errors.ErrCodeInternalServerError)
	}

	return refreshToken, nil
}
