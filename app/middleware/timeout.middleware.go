package middleware

import (
	"gin-boilerplate/package/errors"
	"gin-boilerplate/utils"
	"net/http"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
)

type timeoutMiddleware struct {
}

func NewTimeoutMiddleware() Middleware {
	return &timeoutMiddleware{}
}

func (m *timeoutMiddleware) Handler() gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(60*time.Second),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
		timeout.WithResponse(func(c *gin.Context) {
			err := errors.New(errors.ErrCodeTimeout)
			c.JSON(http.StatusRequestTimeout, utils.FormatErrorResponse(err))
		}),
	)
}
