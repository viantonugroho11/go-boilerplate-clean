package bootstrap

import (
	"strconv"

	"go-boilerplate-clean/internal/config"
	redisinfra "go-boilerplate-clean/internal/infrastructure/cache/redis"

	"github.com/redis/go-redis/v9"
)

func initRedis(cfg *config.Configuration) (*redis.Client, error) {
	return redisinfra.NewClient(cfg.Redis.Addr, cfg.Redis.Password, strconv.Itoa(cfg.Redis.DB))
}
