package global

import (
	"github.com/redis/go-redis/v9"
	db "hieupc05.github/backend-server/db/sqlc"
	"hieupc05.github/backend-server/internal/utils/room"
	"hieupc05.github/backend-server/internal/utils/token"
	"hieupc05.github/backend-server/pkg/logger"
	"hieupc05.github/backend-server/setting"
)

var (
	Config     setting.Config
	Logger     *logger.LoggerZap
	Rdb        *redis.Client
	PgDb       db.Store
	TokenMaker token.Maker
	RoomManage room.Manager
)
