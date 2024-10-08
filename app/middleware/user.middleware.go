package middleware

import (
	_errors "errors"
	"net/http"
	"oauth-server/config"
	"oauth-server/package/errors"
	"oauth-server/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type userMiddleware struct {
}

func NewUserMiddleware() Middleware {
	return &userMiddleware{}
}
func (m *userMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		conf := config.GetConfiguration().Jwt
		resErr := errors.New(errors.ErrCodeUnauthorized)

		token := c.GetHeader(utils.USER_AUTHORIZATION)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.FormatErrorResponse(resErr))
			return
		}
		tokenArr := strings.Split(token, " ")
		if len(tokenArr) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.FormatErrorResponse(resErr))
			return
		}
		if tokenArr[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.FormatErrorResponse(resErr))
			return
		}

		payload, err := utils.VerifyUserToken(tokenArr[1], conf.UserAccessTokenKey)
		if err != nil {
			if _errors.Is(err, jwt.ErrTokenExpired) {
				resErr = errors.New(errors.ErrCodeAccessTokenExpired)
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.FormatErrorResponse(resErr))
			return
		}

		c.Set(string(utils.USER_CONTEXT_KEY), payload)
		c.Next()
	}
}
