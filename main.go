package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"oauth-server/app/controller"
	"oauth-server/app/helper"
	"oauth-server/app/middleware"
	postgres_repository "oauth-server/app/repository/postgres"
	"oauth-server/app/service"
	"oauth-server/config"
	"oauth-server/external/server"
	"oauth-server/package/database"
	logger "oauth-server/package/log"
	_validator "oauth-server/package/validator"
	"oauth-server/utils"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/jutimi/grpc-service/oauth"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
)

func main() {
	conf := config.GetConfiguration()

	// Init uptrace
	uptrace.ConfigureOpentelemetry(
		uptrace.WithDSN(conf.Server.UptraceDNS),
		uptrace.WithServiceName(conf.Server.ServiceName),
		uptrace.WithServiceVersion("1.0.0"),
	)
	tracer := otel.Tracer(conf.Server.ServiceName)

	// Register repositories
	postgresDB := database.GetPostgres()
	// mysqlRepo := mysql_repository.RegisterMysqlRepositories(db)
	postgresRepo := postgres_repository.RegisterPostgresRepositories(postgresDB)

	// Register Others
	helpers := helper.RegisterHelpers(postgresRepo)
	services := service.RegisterServices(helpers, postgresRepo)
	middleware := middleware.RegisterMiddleware()

	// Run GRPC Server
	go startGRPCServer(conf, postgresRepo, helpers)

	// Run gin server
	gin.SetMode(conf.Server.Mode)
	router := gin.Default()
	router.Use(gin.LoggerWithWriter(logger.GetLogger().Writer()))
	router.Use(otelgin.Middleware(conf.Server.ServiceName))

	// Register validator
	v := binding.Validator.Engine().(*validator.Validate)
	v.SetTagName("validate")
	_validator.RegisterCustomValidators(v)

	// Register controllers
	router.GET("/health-check", func(c *gin.Context) {
		c.String(200, "OK")
	})
	controller.RegisterControllers(router, tracer, middleware, services)

	// Start server
	srvErr := make(chan error, 1)
	quit := make(chan os.Signal, 1)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Server.Port),
		Handler: router,
	}
	go func() {
		srvErr <- srv.ListenAndServe()
	}()

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
		uptrace.Shutdown(ctx)
	}()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-srvErr:
		// Error when starting HTTP server.
		return
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}

func init() {
	rootDir := utils.RootDir()
	configFile := fmt.Sprintf("%s/config.yml", rootDir)

	// Init config
	config.Init(configFile)
	// Init database
	database.InitPostgres()
	// Init logger
	logger.Init()

	// Init sentry
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:           config.GetConfiguration().Server.SentryUrl,
		EnableTracing: true,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	}); err != nil {
		log.Fatalf("Error Init Sentry: %s", err.Error())
	}
}

func startGRPCServer(
	conf *config.Configuration,

	postgresRepo postgres_repository.PostgresRepositoryCollections,
	helpers helper.HelperCollections,
) {
	lis, err := net.Listen("tcp", conf.GRPC.OAuthUrl)
	if err != nil {
		log.Fatalf("Error Init GRPC Port: %s", err.Error())
	}
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	// Register server
	oauth.RegisterOAuthRouteServer(grpcServer, server.NewGRPCServer(postgresRepo, helpers))
	oauth.RegisterUserRouteServer(grpcServer, server.NewGRPCServer(postgresRepo, helpers))

	log.Println("Init GRPC Success!")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error Init GRPC: %s", err.Error())
	}
}
