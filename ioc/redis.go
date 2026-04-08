package ioc

import (
	"context"

	"github.com/Knowpals/Knowpals-be-go/config"
	"github.com/redis/go-redis/v9"
)

func InitRedis(conf *config.Config) *redis.Client {
	cmd := redis.NewClient(&redis.Options{Password: conf.Redis.Password, Addr: conf.Redis.Addr, DB: conf.Redis.DB})
	_, err := cmd.Ping(context.Background()).Result()
	if err != nil {
		panic("Redis 连接失败：" + err.Error())
	}
	return cmd
}
