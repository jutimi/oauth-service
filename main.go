package main

import "github.com/gin-gonic/gin"

func main() {
	gin.ForceConsoleColor()

	router := gin.Default()

	router.GET("/health-check", func(c *gin.Context) {
		c.String(200, "OK")
	})

	router.Run(":8080")
}
