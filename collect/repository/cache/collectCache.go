package cache

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"lifelog-grpc/collect/domain"
	"strconv"
	"time"
)

type CollectCache interface {
	Get(ctx context.Context, userId int64) ([]domain.CollectDomain, error)
	Set(ctx context.Context, userId int64, cds []domain.CollectDomain) error
}

type CollectRedisCache struct {
	cmd redis.Cmdable
}

func NewCollectRedisCache(cmd redis.Cmdable) CollectCache {
	return &CollectRedisCache{
		cmd: cmd,
	}
}

func (c *CollectRedisCache) Get(ctx context.Context, userId int64) ([]domain.CollectDomain, error) {
	// 获取key
	key := c.getKey(userId)
	result, err := c.cmd.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	var ads []domain.CollectDomain
	// 反序列化
	err = json.Unmarshal(result, &ads)
	if err != nil {
		return nil, err
	}
	return ads, nil
}

func (c *CollectRedisCache) Set(ctx context.Context, userId int64, cds []domain.CollectDomain) error {
	// 获取key
	key := c.getKey(userId)
	// 序列化
	marshal, err := json.Marshal(cds)
	if err != nil {
		return err
	}
	err = c.cmd.Set(ctx, key, marshal, 5*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *CollectRedisCache) getKey(userId int64) string {
	return "collect_userId_" + strconv.Itoa(int(userId))
}
