package bootstrap

import (
	"strconv"

	redisinfra "go-boilerplate-clean/internal/infrastructure/cache/redis"

	"github.com/redis/go-redis/v9"
)

// InitRedis buat Redis client. Pakai Config() global.
func initRedis() (*redis.Client, error) {
	c := Config()
	return redisinfra.NewClient(c.Redis.Addr, c.Redis.Password, strconv.Itoa(c.Redis.DB))
}
