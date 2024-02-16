package main

import (
	"fmt"
	"gin-boilerplate/app/controller"
	"gin-boilerplate/config"
	"gin-boilerplate/package/database"
	logger "gin-boilerplate/package/log"
	_validator "gin-boilerplate/package/validator"
	"gin-boilerplate/utils"

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

	// Register controllers
	router.GET("/health-check", func(c *gin.Context) {
		c.String(200, "OK")
	})
	registerControllers(router)

	// Register services

	router.Run(fmt.Sprintf(":%d", config.GetConfiguration().Server.Port))
}

func init() {
	rootDir := utils.RootDir()
	configFile := fmt.Sprintf("%s/config.yml", rootDir)

	config.Init(configFile)
	database.Init()
	logger.Init()
}

func registerControllers(router *gin.Engine) {
	controller.NewAPIUserController(router)
}
