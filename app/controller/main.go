package controller

import (
	"oauth-server/app/controller/api"
	"oauth-server/app/middleware"
	"oauth-server/app/service"

	"github.com/gin-gonic/gin"
)

func RegisterControllers(
	router *gin.Engine,
	services service.ServiceCollections,
	middleware middleware.MiddlewareCollections,
) {
	api.NewApiUserController(router, services)
	api.NewApiOAuthController(router, services, middleware)
}
