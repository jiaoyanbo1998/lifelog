package cache

import (
	"errors"
	"fmt"
	"github.com/patrickmn/go-cache"
	"lifelog-grpc/lifeLog/domain"
	"time"
)

// LocalCache 本地缓存
type LocalCache struct {
	cache *cache.Cache // 使用 go-cache 的缓存
}

// NewLocalCache 创建一个新的LocalCache
func NewLocalCache(cache *cache.Cache) *LocalCache {
	// 创建一个缓存实例，设置默认过期时间和清理间隔
	return &LocalCache{
		cache: cache,
	}
}

// Set 设置详情页缓存
func (c *LocalCache) Set(lifeLogDomain domain.LifeLogDomain) error {
	// 存储数据副本防止外部修改影响缓存
	var log domain.LifeLogDomain
	log = lifeLogDomain
	key := fmt.Sprintf("Detail_%d", lifeLogDomain.Author.Id)
	c.cache.Set(key, log, 1*time.Minute)
	return nil
}

// Get 获取详情页缓存
func (c *LocalCache) Get(authorId int64) (domain.LifeLogDomain, error) {
	// 从本地缓存中获取数据
	key := fmt.Sprintf("Detail_%d", authorId)
	value, found := c.cache.Get(key)
	if !found {
		return domain.LifeLogDomain{}, errors.New("not found 本地缓存")
	}
	// 类型断言，确保数据类型正确
	log, ok := value.(domain.LifeLogDomain)
	if !ok {
		return domain.LifeLogDomain{}, errors.New("非法数据类型")
	}
	// 返回数据副本防止外部修改影响缓存
	return log, nil
}

// GetFirstPage 获取本地缓存
func (c *LocalCache) GetFirstPage(authorId int64) ([]domain.LifeLogDomain, error) {
	// 从本地缓存中获取数据
	key := fmt.Sprintf("FirstPage_%d", authorId)
	value, found := c.cache.Get(key)
	if !found {
		return nil, errors.New("not found 本地缓存")
	}
	// 类型断言，确保数据类型正确
	logs, ok := value.([]domain.LifeLogDomain)
	if !ok {
		return nil, errors.New("非法数据类型")
	}
	// 返回数据副本防止外部修改影响缓存
	copied := make([]domain.LifeLogDomain, len(logs))
	copy(copied, logs)
	return copied, nil
}

// SetFirstPage 设置本地缓存
func (c *LocalCache) SetFirstPage(authorId int64, lifeLogDomains []domain.LifeLogDomain) error {
	// 存储数据副本防止外部修改影响缓存
	logs := make([]domain.LifeLogDomain, len(lifeLogDomains))
	copy(logs, lifeLogDomains)
	key := fmt.Sprintf("FirstPage_%d", authorId)
	c.cache.Set(key, logs, 1*time.Minute)
	return nil
}

// DelFirstPage 删除本地缓存
func (c *LocalCache) DelFirstPage(authorId int64) error {
	key := fmt.Sprintf("FirstPage_%d", authorId)
	c.cache.Delete(key)
	return nil
}
