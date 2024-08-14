package api

import (
	"context"
	"net/http"
	"oauth-server/app/middleware"
	"oauth-server/app/model"
	"oauth-server/app/service"
	_errors "oauth-server/package/errors"
	"oauth-server/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.opentelemetry.io/otel/trace"
)

type oAuthHandler struct {
	tracer     trace.Tracer
	middleware middleware.MiddlewareCollections
	services   service.ServiceCollections
}

func NewApiOAuthController(
	router *gin.Engine,
	tracer trace.Tracer,
	middleware middleware.MiddlewareCollections,
	services service.ServiceCollections,
) {
	handler := oAuthHandler{tracer, middleware, services}

	group := router.Group("api/v1/oauth")
	{
		userGroup := group.Group("users")
		userGroup.POST("/refresh", handler.refreshUserToken)
		userGroup.POST("/login", handler.userLogin)
		userGroup.POST("/logout", middleware.UserMW.Handler(), handler.userLogout)

		WorkspaceGroup := userGroup.Group("workspaces")
		WorkspaceGroup.POST("/refresh", handler.refreshWorkspaceToken)
		WorkspaceGroup.POST("/login", middleware.UserMW.Handler(), handler.WorkspaceLogin)
		WorkspaceGroup.POST("/logout", middleware.WorkspaceMW.Handler(), handler.WorkspaceLogout)
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
	ctx = context.WithValue(ctx, utils.GIN_CONTEXT_KEY, c)
	ctx = context.WithValue(ctx, utils.SCOPE_CONTEXT_KEY, utils.USER_SCOPE)
	ctx, main := h.tracer.Start(ctx, "refresh-user-token")
	defer func() {
		cancel()
		main.End()
	}()

	res, err := h.services.OAuthSvc.RefreshToken(ctx, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(res))
}

func (h *oAuthHandler) refreshWorkspaceToken(c *gin.Context) {
	var data model.RefreshTokenRequest
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		resErr := _errors.NewValidatorError(err)
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(resErr))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	ctx = context.WithValue(ctx, utils.GIN_CONTEXT_KEY, c)
	ctx = context.WithValue(ctx, utils.SCOPE_CONTEXT_KEY, utils.WORKSPACE_SCOPE)
	ctx, main := h.tracer.Start(ctx, "refresh-Workspace-token")
	defer func() {
		cancel()
		main.End()
	}()

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
	ctx = context.WithValue(ctx, utils.GIN_CONTEXT_KEY, c)
	ctx = context.WithValue(ctx, utils.SCOPE_CONTEXT_KEY, utils.USER_SCOPE)
	ctx, main := h.tracer.Start(ctx, "user-login")
	defer func() {
		cancel()
		main.End()
	}()

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
	ctx = context.WithValue(ctx, utils.GIN_CONTEXT_KEY, c)
	ctx = context.WithValue(ctx, utils.SCOPE_CONTEXT_KEY, utils.USER_SCOPE)
	ctx, main := h.tracer.Start(ctx, "user-logout")
	defer func() {
		cancel()
		main.End()
	}()

	res, err := h.services.OAuthSvc.Logout(ctx, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(res))
}

func (h *oAuthHandler) WorkspaceLogin(c *gin.Context) {
	var data model.WorkspaceLoginRequest
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		resErr := _errors.NewValidatorError(err)
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(resErr))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	ctx = context.WithValue(ctx, utils.GIN_CONTEXT_KEY, c)
	ctx = context.WithValue(ctx, utils.SCOPE_CONTEXT_KEY, utils.WORKSPACE_SCOPE)
	ctx, main := h.tracer.Start(ctx, "Workspace-login")
	defer func() {
		cancel()
		main.End()
	}()

	res, err := h.services.OAuthSvc.Login(ctx, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(res))
}

func (h *oAuthHandler) WorkspaceLogout(c *gin.Context) {
	var data model.UserLogoutRequest
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		resErr := _errors.NewValidatorError(err)
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(resErr))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	ctx = context.WithValue(ctx, utils.GIN_CONTEXT_KEY, c)
	ctx = context.WithValue(ctx, utils.SCOPE_CONTEXT_KEY, utils.WORKSPACE_SCOPE)
	ctx, main := h.tracer.Start(ctx, "Workspace-logout")
	defer func() {
		cancel()
		main.End()
	}()

	res, err := h.services.OAuthSvc.Logout(ctx, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(res))
}
