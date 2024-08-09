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
	"go.opentelemetry.io/otel/trace"
)

type userHandler struct {
	tracer   trace.Tracer
	services service.ServiceCollections
}

func NewApiUserController(
	router *gin.Engine,
	tracer trace.Tracer,
	services service.ServiceCollections,
) {
	handler := userHandler{tracer, services}

	group := router.Group("api/v1/users")
	{
		group.POST("/register", handler.register)
	}
}

func (h *userHandler) register(c *gin.Context) {
	var data model.RegisterRequest
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		resErr := _errors.NewValidatorError(err)
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(resErr))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	ctx = context.WithValue(ctx, utils.GIN_CONTEXT_KEY, c)
	ctx, main := h.tracer.Start(ctx, "register")
	defer func() {
		cancel()
		main.End()
	}()

	res, err := h.services.UserSvc.Register(ctx, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(res))
}
