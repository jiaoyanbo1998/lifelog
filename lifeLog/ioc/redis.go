package ioc

import (
	"github.com/redis/go-redis/v9"
)

//var (
//	redisClient redis.Cmdable
//	once        sync.Once
//)
//
//func GetRedis() redis.Cmdable {
//	once.Do(func() {
//		redisClient = initRedis()
//	})
//	return redisClient
//}

func InitRedis() redis.Cmdable {
	// 初始化Redis客户端
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // 如果有密码则设置
		DB:       0,  // 使用默认数据库
	})
	return client
}
