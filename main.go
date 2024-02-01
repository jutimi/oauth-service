package main

import (
	"fmt"
	"gin-boilerplate/config"
	"gin-boilerplate/package/database"
	logger "gin-boilerplate/package/log"

	"github.com/gin-gonic/gin"
)

func main() {
	Init()

	router := gin.Default()
	router.Use(gin.LoggerWithWriter(logger.GetLogger().Writer()))

	router.GET("/health-check", func(c *gin.Context) {
		c.String(200, "OK")
	})

	router.Run(fmt.Sprintf(":%d", config.GetConfiguration().Server.Port))
}

func Init() {
	config.Init("config.yml")
	database.Init()
	logger.Init()
}
