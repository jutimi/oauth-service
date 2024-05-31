package controller

import (
	"gin-boilerplate/app/controller/api"
	"gin-boilerplate/app/service"

	"github.com/gin-gonic/gin"
)

func RegisterControllers(router *gin.Engine, services service.ServiceCollections) {
	api.NewApiUserController(router, services)
	api.NewApiOAuthController(router, services)
}
