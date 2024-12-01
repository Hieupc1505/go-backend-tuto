package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"hieupc05.github/backend-server/internal/services"
)

type UserController struct {
	getInfo *services.Services
}

func NewUserController() *UserController {
	return &UserController{
		getInfo: services.NewService(),
	}
}

func (ur *UserController) User(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "User",
		"user":    ur.getInfo.GetUserInfoService(),
	})
}
