package main

import (
	"context"
	"fmt"
	"gin-boilerplate/app/controller"
	"gin-boilerplate/app/helper"
	repository "gin-boilerplate/app/repository/mysql"
	"gin-boilerplate/app/service"
	"gin-boilerplate/config"
	"gin-boilerplate/package/database"
	logger "gin-boilerplate/package/log"
	_validator "gin-boilerplate/package/validator"
	"gin-boilerplate/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func main() {
	router := gin.Default()
	router.Use(gin.LoggerWithWriter(logger.GetLogger().Writer()))

	// Register validator
	v := binding.Validator.Engine().(*validator.Validate)
	_validator.RegisterCustomValidators(v)

	// Register repositories
	db := database.GetPostgres()
	mysqlRepo := repository.RegisterMysqlRepositories(db)
	helpers := helper.RegisterHelpers(mysqlRepo)
	services := service.RegisterServices(helpers, mysqlRepo)

	// Register controllers
	router.GET("/health-check", func(c *gin.Context) {
		c.String(200, "OK")
	})
	controller.RegisterControllers(router, services)

	// Start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.GetConfiguration().Server.Port),
		Handler: router,
	}

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
}

func init() {
	rootDir := utils.RootDir()
	configFile := fmt.Sprintf("%s/config.yml", rootDir)

	config.Init(configFile)
	database.InitPostgres()
	logger.Init()
}
