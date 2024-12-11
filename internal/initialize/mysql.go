package initialize

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	db "hieupc05.github/backend-server/db/sqlc"
	"hieupc05.github/backend-server/global"
)

func InitMysql() {
	connPool, err := pgxpool.New(context.Background(), global.Config.PgDb.DbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	global.PgDb = db.NewStore(connPool)
}
