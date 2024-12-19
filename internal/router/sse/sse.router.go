package sse

import (
	"github.com/gin-gonic/gin"
	"hieupc05.github/backend-server/global"
	controllers "hieupc05.github/backend-server/internal/controller"
	"hieupc05.github/backend-server/internal/middlewares"
)

type SseMessage struct{}

func (s *SseMessage) InitSseRouter(Router *gin.RouterGroup) {
	Router.Use(middlewares.AuthenMiddleware(global.TokenMaker))
	sseMessageRouter := Router.Group("/sse")
	sseManagerContest := Router.Group("/contest")
	sseCtrl := controllers.NewSseController()
	{
		sseMessageRouter.GET("", sseCtrl.SseConnection)
	}
	{
		sseManagerContest.GET("/start", sseCtrl.SseStartContest)
	}
}
