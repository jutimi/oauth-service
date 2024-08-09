package controller

import (
	"oauth-server/app/controller/api"
	"oauth-server/app/controller/cms"
	"oauth-server/app/middleware"
	"oauth-server/app/service"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

func RegisterControllers(
	router *gin.Engine,
	tracer trace.Tracer,
	middleware middleware.MiddlewareCollections,
	services service.ServiceCollections,
) {
	api.NewApiUserController(router, tracer, services)
	api.NewApiOAuthController(router, tracer, middleware, services)

	cms.NewApiPermissionController(router, services)
}
