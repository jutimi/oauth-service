package main

import (
	"fmt"
	"gin-boilerplate/config"
	"gin-boilerplate/package/database"

	"github.com/gin-gonic/gin"
)

func main() {
	Init()

	gin.ForceConsoleColor()

	router := gin.Default()

	router.GET("/health-check", func(c *gin.Context) {
		c.String(200, "OK")
	})

	router.Run(fmt.Sprintf(":%d", config.GetConfiguration().Server.Port))
}

func Init() {
	config.Init("config.yml")
	database.Init()
}
