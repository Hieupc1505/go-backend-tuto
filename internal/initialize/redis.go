package initialize

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"hieupc05.github/backend-server/global"
)

func InitRedis() {
	m := global.Config.Redis

	var s = fmt.Sprintf("%s:%d", m.Host_name, m.Port)
	db := newClient(s, m.Password)

	global.Rdb = db

}

func newClient(addr string, password string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})
}
