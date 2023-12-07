package ioc

import (
	"github.com/redis/go-redis/v9"
	"github.com/xuqil/webook/config"
)

func InitRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	return redisClient
}
