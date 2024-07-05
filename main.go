package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"oauth-server/app/controller"
	"oauth-server/app/helper"
	postgres_repository "oauth-server/app/repository/postgres"
	"oauth-server/app/service"
	"oauth-server/config"
	user_grpc "oauth-server/gRPC/user"
	"oauth-server/package/database"
	logger "oauth-server/package/log"
	_validator "oauth-server/package/validator"
	"oauth-server/utils"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
)

func main() {
	conf := config.GetConfiguration()

	// Register repositories
	postgresDB := database.GetPostgres()
	// mysqlRepo := mysql_repository.RegisterMysqlRepositories(db)
	postgresRepo := postgres_repository.RegisterPostgresRepositories(postgresDB)

	// Register Others
	helpers := helper.RegisterHelpers(postgresRepo)
	services := service.RegisterServices(helpers, postgresRepo)

	// Run GRPC Server
	go startGRPCServer(conf, postgresRepo)

	// Run gin server
	gin.SetMode(conf.Server.Mode)
	router := gin.Default()
	router.Use(gin.LoggerWithWriter(logger.GetLogger().Writer()))

	// Register validator
	v := binding.Validator.Engine().(*validator.Validate)
	v.SetTagName("validate")
	_validator.RegisterCustomValidators(v)

	// Register controllers
	router.GET("/health-check", func(c *gin.Context) {
		c.String(200, "OK")
	})
	controller.RegisterControllers(router, services)

	// Start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Server.Port),
		Handler: router,
	}

	if !gin.IsDebugging() {
		go func() {
			// service connections
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutdown Server ...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}

		// catching ctx.Done(). timeout of 5 seconds.
		select {
		case <-ctx.Done():
			log.Println("timeout of 5 seconds.")
		}
		log.Println("Server exiting")
	} else {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}
}

func init() {
	rootDir := utils.RootDir()
	configFile := fmt.Sprintf("%s/config.yml", rootDir)

	config.Init(configFile)
	database.InitPostgres()
	logger.Init()
}

func startGRPCServer(
	conf *config.Configuration,

	postgresRepo postgres_repository.PostgresRepositoryCollections,
) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.GRPC.Port))
	if err != nil {
		panic(err)
	}
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	// Register serever
	user_grpc.RegisterUserRouteServer(grpcServer, user_grpc.NewUserServer(postgresRepo))

	logger.Println(logger.LogPrintln{
		FileName:  "main.go",
		FuncName:  "main",
		TraceData: "",
		Msg: fmt.Sprintf(
			"GRPC Server Running Port - %d",
			conf.GRPC.Port,
		),
	})
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
