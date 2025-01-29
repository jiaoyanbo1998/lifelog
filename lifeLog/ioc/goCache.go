package ioc

import (
	"github.com/patrickmn/go-cache"
	"time"
)

func InitGoCache() *cache.Cache {
	// 参数1：过期时间
	// 参数2：清理间隔（每隔多长时间清理一次过期缓存）
	return cache.New(1*time.Minute, 10*time.Minute)
}
