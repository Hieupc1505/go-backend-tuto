package routers

import (
	"hieupc05.github/backend-server/internal/router/manager"
	"hieupc05.github/backend-server/internal/router/sse"
	"hieupc05.github/backend-server/internal/router/upload"
	"hieupc05.github/backend-server/internal/router/user"
)

type RouterGroup struct {
	User    user.UserRouterGroup
	Manager manager.ManagerRouterGroup
	Upload  upload.UploadRouterGroup
	Sse     sse.SseRouterGroup
}

var RouterGroupApp = new(RouterGroup)
