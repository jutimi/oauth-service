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

type accountHandler struct {
	services service.ServiceCollections
}

func NewAPIUserController(router *gin.Engine, services service.ServiceCollections) {
	handler := accountHandler{services}

	group := router.Group("api/v1/users")
	{
		group.POST("/login", handler.login)
		group.POST("/register", handler.register)
		group.POST("/logout", handler.logout)
	}
}

func (h *accountHandler) login(c *gin.Context) {
	var data model.LoginRequest
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		resErr := _errors.NewValidatorError(err)
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(resErr))

		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	ctx = context.WithValue(ctx, utils.GIN_CONTEXT_KEY, c)

	res, err := h.services.UserSvc.Login(ctx, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(err))
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(res))
}

func (h *accountHandler) register(c *gin.Context) {
	var data model.RegisterRequest
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		resErr := _errors.NewValidatorError(err)
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(resErr))

		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	ctx = context.WithValue(ctx, "gin", c)

	res, err := h.services.UserSvc.Register(ctx, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(err))
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(res))
}

func (h *accountHandler) logout(c *gin.Context) {
	var data model.LogoutRequest
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		resErr := _errors.NewValidatorError(err)
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(resErr))

		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	ctx = context.WithValue(ctx, "gin", c)

	res, err := h.services.UserSvc.Logout(ctx, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(err))
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(res))
}
