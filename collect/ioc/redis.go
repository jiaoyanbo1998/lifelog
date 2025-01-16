package ioc

import (
	"github.com/redis/go-redis/v9"
	"sync"
)

var (
	redisClient redis.Cmdable
	redisOnce   sync.Once
)

// GetRedis 懒汉式单例模式
func GetRedis() redis.Cmdable {
	redisOnce.Do(func() {
		redisClient = InitRedis()
	})
	return redisClient
}

func InitRedis() redis.Cmdable {
	// 初始化Redis客户端
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // 如果有密码则设置
		DB:       0,  // 使用默认数据库
	})
	return client
}
