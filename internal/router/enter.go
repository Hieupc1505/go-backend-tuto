package routers

import (
	"hieupc05.github/backend-server/internal/router/manager"
	"hieupc05.github/backend-server/internal/router/upload"
	"hieupc05.github/backend-server/internal/router/user"
)

type RouterGroup struct {
	User    user.UserRouterGroup
	Manager manager.ManagerRouterGroup
	Upload  upload.UploadRouterGroup
}

var RouterGroupApp = new(RouterGroup)
