package routers

import (
	"github.com/gin-gonic/gin"
	c "hieupc05.github/backend-server/internal/controller"
	"hieupc05.github/backend-server/internal/middlewares"
)

func NewRoute() *gin.Engine {
	r := gin.Default()
	//use the midddlewarre
	r.Use(middlewares.AuthenMiddleware())
	v1 := r.Group("/v1/2024")
	{
		v1.GET("/ping", c.NewPongController().Pong)
		v1.GET("/user", c.NewUserController().User)

	}

	return r
}
