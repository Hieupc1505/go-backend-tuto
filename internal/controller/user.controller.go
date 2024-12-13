package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"hieupc05.github/backend-server/global"
	"hieupc05.github/backend-server/internal/services"
	"hieupc05.github/backend-server/response"
)

type RegisterForm struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type VerifyOTPForm struct {
	OTP   int    `json:"otp" binding:"required,min=100000,max=999999"`
	Email string `json:"email" binding:"required,email"`
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserController struct {
	userService services.IUserServices
}

func NewUserController(userService services.IUserServices) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (uc *UserController) Register(c *gin.Context) {
	var req RegisterForm
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error(err.Error(), zap.Error(err))
		res := response.ErrorResponse(response.ErrAuthFail)
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// Gọi service
	result, httpStatus := uc.userService.Register(req.Email, req.Password)
	// result, httpStatus := response.Response{}, 200

	// Trả về HTTP response tương ứng
	c.JSON(httpStatus, result)
}

func (uc *UserController) CreateUser(c *gin.Context) {
	var req VerifyOTPForm
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("error body")
		global.Logger.Error(err.Error(), zap.Error(err))
		res := response.ErrorResponse(response.ErrAuthFail)
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// Gọi service
	result, httpStatus := uc.userService.CreateUser(c, req.OTP, req.Email)

	// Trả về HTTP response tương ứng
	c.JSON(httpStatus, result)
}

func (uc *UserController) Login(c *gin.Context) {
	var req loginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error(err.Error(), zap.Error(err))
		res := response.ErrorResponse(response.ErrAuthFail)
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// Gọi service
	result, httpStatus := uc.userService.Login(c, req.Email, req.Password)

	// Trả về HTTP response tương ứng
	c.JSON(httpStatus, result)
}

func (uc *UserController) GetUserInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"userName": "hieupc05",
		"email":    "hieupc05@hieupc05.com",
	})
}
