package initialize

import (
	"hieupc05.github/backend-server/global"
	"hieupc05.github/backend-server/pkg/logger"
)

func InitLogger() {
	global.Logger = logger.NewLogger(global.Config.Logger)
}
