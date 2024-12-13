package initialize

import (
	"github.com/gin-gonic/gin"
	"hieupc05.github/backend-server/global"
	routers "hieupc05.github/backend-server/internal/router"
)

func InitRouter() *gin.Engine {
	var r *gin.Engine
	if global.Config.Server.Mode == "dev" {
		gin.SetMode(gin.DebugMode)
		gin.ForceConsoleColor()
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
	}

	//middleware
	// r.Use()// logging
	// r.Use() // cross
	// r.Use() // limiter global

	manageRoute := routers.RouterGroupApp.Manager
	userRoute := routers.RouterGroupApp.User
	uploadRoute := routers.RouterGroupApp.Upload
	sseRoute := routers.RouterGroupApp.Sse
	contestRoute := routers.RouterGroupApp.Contest

	MainGroup := r.Group("/v1/2024")
	{
		MainGroup.GET("/checkStatus")
	}
	{
		userRoute.InitUserRouter(MainGroup)
		userRoute.InitProductRouter(MainGroup)
	}
	{
		manageRoute.InitAdminRouter(MainGroup)
		manageRoute.InitUserManagerRouter(MainGroup)
	}
	{
		uploadRoute.InitUploadImageRouter(MainGroup)
	}
	{
		sseRoute.InitSseRouter(MainGroup)
	}
	{
		contestRoute.InitContestRouter(MainGroup)
	}

	return r
}
