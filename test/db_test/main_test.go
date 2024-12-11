package db_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	db "hieupc05.github/backend-server/db/sqlc"
	"hieupc05.github/backend-server/global"
	"hieupc05.github/backend-server/internal/initialize"
)

var testStore *db.Queries

func TestMain(m *testing.M) {
	initialize.LoadConfig("../../configs/")
	connPool, err := pgxpool.New(context.Background(), global.Config.PgDb.DbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testStore = db.New(connPool)
	os.Exit(m.Run())
}
