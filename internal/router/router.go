package routers

import (
	"github.com/gin-gonic/gin"
)

func NewRoute() *gin.Engine {
	r := gin.Default()
	//use the midddlewarre
	// r.Use(middlewares.AuthenMiddleware())
	v1 := r.Group("/v1/2024")
	{
		v1.GET("/ping")
		v1.GET("/user")

	}

	return r
}
