package controller

import (
	"gin-boilerplate/app/model"
	_errors "gin-boilerplate/package/errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func NewAPIUserController(router *gin.Engine) {
	group := router.Group("api/v1/users")
	{
		group.POST("/login", login)
		group.POST("/register", register)
		group.POST("/logout", logout)
	}
}

func login(c *gin.Context) {
	var data model.LoginRequest
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		_errors.HandleValidatorError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

func register(c *gin.Context) {

}

func logout(c *gin.Context) {

}
