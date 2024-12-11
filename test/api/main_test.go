package api

import (
	"github.com/gin-gonic/gin"
	db "hieupc05.github/backend-server/db/sqlc"
	"hieupc05.github/backend-server/global"
	"hieupc05.github/backend-server/internal/initialize"
)

func newTestServer(store db.Store) *gin.Engine {
	initialize.LoadConfig("../../configs/")
	initialize.InitLogger()
	initialize.InitRedis()
	global.PgDb = store
	router := initialize.InitRouter()
	return router
}
