package initialize

import (
	"go.uber.org/zap"
	"hieupc05.github/backend-server/global"
)

func Run() {
	//load config
	LoadConfig("./configs/")
	InitLogger()
	global.Logger.Info("Config Log ok!!", zap.String("OK", "Success"))

	InitMysql()
	InitRedis()
	InitCommon()

	r := InitRouter()
	r.Run(":8002")

}
