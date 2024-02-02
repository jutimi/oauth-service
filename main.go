package main

import (
	"fmt"
	"gin-boilerplate/config"
	"gin-boilerplate/package/database"
	logger "gin-boilerplate/package/log"
	"gin-boilerplate/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(gin.LoggerWithWriter(logger.GetLogger().Writer()))

	router.GET("/health-check", func(c *gin.Context) {
		c.String(200, "OK")
	})

	router.Run(fmt.Sprintf(":%d", config.GetConfiguration().Server.Port))
}

func init() {
	rootDir := utils.RootDir()
	configFile := fmt.Sprintf("%s/config.yml", rootDir)

	config.Init(configFile)
	database.Init()
	logger.Init()
}
