package helper

import (
	"fmt"
	"gin-boilerplate/app/entity"
	"gin-boilerplate/config"
	logger "gin-boilerplate/package/log"
	"gin-boilerplate/utils"
)

type oauthHelper struct {
}

func NewOauthHelper() OauthHelper {
	return &oauthHelper{}
}

func (h *oauthHelper) GenerateAccessToken(user entity.User) (string, error) {
	conf := config.GetConfiguration().Jwt

	payload := &utils.UserPayload{
		ID: user.ID,
	}
	accessToken, err := utils.GenerateToken(payload, conf.UserAccessTokenKey, utils.USER_ACCESS_TOKEN_IAT)
	if err != nil {
		logger.Println(logger.LogPrintln{
			FileName:  "app/helper/oauth.helper.go",
			FuncName:  "GenerateAccessToken",
			TraceData: fmt.Sprintf("%s/%s", *user.Email, *user.PhoneNumber),
			Msg:       err.Error(),
		})
		return "", err
	}

	return accessToken, nil
}

func (h *oauthHelper) GenerateRefreshToken(user entity.User) (string, error) {
	conf := config.GetConfiguration().Jwt

	payload := &utils.UserPayload{
		ID: user.ID,
	}
	refreshToken, err := utils.GenerateToken(payload, conf.UserRefreshTokenKey, utils.USER_REFRESH_TOKEN_IAT)
	if err != nil {
		logger.Println(logger.LogPrintln{
			FileName:  "app/service/user.service.go",
			FuncName:  "GenerateRefreshToken",
			TraceData: fmt.Sprintf("%s/%s", *user.Email, *user.PhoneNumber),
			Msg:       err.Error(),
		})
		return "", err
	}

	return refreshToken, nil
}
