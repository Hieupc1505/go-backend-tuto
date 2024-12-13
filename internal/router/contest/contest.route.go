package contest

import (
	"github.com/gin-gonic/gin"
	"hieupc05.github/backend-server/global"
	controllers "hieupc05.github/backend-server/internal/controller"
	"hieupc05.github/backend-server/internal/middlewares"
	"hieupc05.github/backend-server/internal/services"
)

type ContestRouter struct{}

func (r *ContestRouter) InitContestRouter(Router *gin.RouterGroup) {

	sv := services.NewContestService()
	contestControler := controllers.NewContestController(sv)

	Router.Use(middlewares.AuthenMiddleware(global.TokenMaker))
	contestRouter := Router.Group("/contest")
	{
		contestRouter.POST("/create", contestControler.CreateContest)
		contestRouter.GET("/live/me")
		contestRouter.GET("/start/:id", contestControler.LiveContest)
		contestRouter.POST("/play/:id", contestControler.PlayContest)
		contestRouter.GET("/live/:id")
		contestRouter.GET("/stop/:id", contestControler.EndContest)
		contestRouter.POST("/:id/submit-paper", contestControler.SubmitAnswer)
	}
}
