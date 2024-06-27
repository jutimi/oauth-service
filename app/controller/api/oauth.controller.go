package api

import (
	"context"
	"gin-boilerplate/app/model"
	"gin-boilerplate/app/service"
	_errors "gin-boilerplate/package/errors"
	"gin-boilerplate/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type oAuthHandler struct {
	services service.ServiceCollections
}

func NewApiOAuthController(router *gin.Engine, services service.ServiceCollections) {
	handler := oAuthHandler{services}

	group := router.Group("api/v1/oauth")
	{
		group.POST("/refresh", handler.refreshToken)
	}
}

func (h *oAuthHandler) refreshToken(c *gin.Context) {
	var data model.RefreshTokenRequest
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		resErr := _errors.NewValidatorError(err)
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(resErr))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, utils.GIN_CONTEXT_KEY, c)

	res, err := h.services.OAuthSvc.RefreshToken(ctx, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(res))
}
