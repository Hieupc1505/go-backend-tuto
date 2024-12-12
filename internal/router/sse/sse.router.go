package sse

import (
	"github.com/gin-gonic/gin"
	controllers "hieupc05.github/backend-server/internal/controller"
)

type SseMessage struct{}

func (s *SseMessage) InitSseRouter(Router *gin.RouterGroup) {
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
