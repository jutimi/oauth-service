package api

import (
	"context"
	"net/http"
	"oauth-server/app/model"
	"oauth-server/app/service"
	_errors "oauth-server/package/errors"
	"oauth-server/utils"
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
		userGroup := group.Group("users")
		userGroup.POST("/refresh", handler.refreshUserToken)
		userGroup.POST("/login", handler.userLogin)
		userGroup.POST("/logout", handler.userLogout)

		wsGroup := userGroup.Group("workspaces")
		wsGroup.POST("/refresh", handler.refreshWSToken)
		wsGroup.POST("/login", handler.wsLogin)
		wsGroup.POST("/logout", handler.wsLogout)
	}
}

func (h *oAuthHandler) refreshUserToken(c *gin.Context) {
	var data model.RefreshTokenRequest
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		resErr := _errors.NewValidatorError(err)
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(resErr))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, utils.GIN_CONTEXT_KEY, c)
	ctx = context.WithValue(ctx, utils.SCOPE_CONTEXT_KEY, utils.USER_SCOPE)

	res, err := h.services.OAuthSvc.RefreshToken(ctx, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(res))
}

func (h *oAuthHandler) refreshWSToken(c *gin.Context) {
	var data model.RefreshTokenRequest
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		resErr := _errors.NewValidatorError(err)
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(resErr))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, utils.GIN_CONTEXT_KEY, c)
	ctx = context.WithValue(ctx, utils.SCOPE_CONTEXT_KEY, utils.WORKSPACE_SCOPE)

	res, err := h.services.OAuthSvc.RefreshToken(ctx, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(res))
}

func (h *oAuthHandler) userLogin(c *gin.Context) {
	var data model.UserLoginRequest
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		resErr := _errors.NewValidatorError(err)
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(resErr))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, utils.GIN_CONTEXT_KEY, c)
	ctx = context.WithValue(ctx, utils.SCOPE_CONTEXT_KEY, utils.USER_SCOPE)

	res, err := h.services.OAuthSvc.Login(ctx, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(res))
}

func (h *oAuthHandler) userLogout(c *gin.Context) {
	var data model.UserLogoutRequest
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		resErr := _errors.NewValidatorError(err)
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(resErr))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, utils.GIN_CONTEXT_KEY, c)
	ctx = context.WithValue(ctx, utils.SCOPE_CONTEXT_KEY, utils.USER_SCOPE)

	res, err := h.services.OAuthSvc.Logout(ctx, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(res))
}

func (h *oAuthHandler) wsLogin(c *gin.Context) {
	var data model.WorkspaceLoginRequest
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		resErr := _errors.NewValidatorError(err)
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(resErr))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, utils.GIN_CONTEXT_KEY, c)
	ctx = context.WithValue(ctx, utils.SCOPE_CONTEXT_KEY, utils.WORKSPACE_SCOPE)

	res, err := h.services.OAuthSvc.Login(ctx, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(res))
}

func (h *oAuthHandler) wsLogout(c *gin.Context) {
	var data model.UserLogoutRequest
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		resErr := _errors.NewValidatorError(err)
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(resErr))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, utils.GIN_CONTEXT_KEY, c)
	ctx = context.WithValue(ctx, utils.SCOPE_CONTEXT_KEY, utils.WORKSPACE_SCOPE)

	res, err := h.services.OAuthSvc.Logout(ctx, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(res))
}
