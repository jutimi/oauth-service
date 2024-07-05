package controller

import (
	"oauth-server/app/controller/api"
	"oauth-server/app/service"

	"github.com/gin-gonic/gin"
)

func RegisterControllers(router *gin.Engine, services service.ServiceCollections) {
	api.NewApiUserController(router, services)
	api.NewApiOAuthController(router, services)
}
