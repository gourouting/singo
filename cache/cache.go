package cache

import (
	"context"
	"os"
	"singo/util"
	"strconv"

	"github.com/redis/go-redis/v9"
)

// RedisClient Redis缓存客户端单例
var RedisClient *redis.Client

// Redis 在中间件中初始化redis链接
func Redis() {
	db, _ := strconv.ParseUint(os.Getenv("REDIS_DB"), 10, 64)
	client := redis.NewClient(&redis.Options{
		Addr:       os.Getenv("REDIS_ADDR"),
		Password:   os.Getenv("REDIS_PW"),
		DB:         int(db),
		MaxRetries: 1,
	})

	_, err := client.Ping(context.Background()).Result()

	if err != nil {
		util.Log().Panic("连接Redis不成功", err)
	}

	RedisClient = client
}
