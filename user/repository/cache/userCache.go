package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type UserCache interface {
	Set(ctx context.Context, sessionId string) error
	SetHistoryPassword(ctx context.Context, userKey, code string) error
	GetHistoryPassword(ctx context.Context, userKey string) ([]string, error)
}

type UserRedisCache struct {
	cmd redis.Cmdable
}

func NewUserCache(cmd redis.Cmdable) UserCache {
	return &UserRedisCache{
		cmd: cmd,
	}
}

func (u *UserRedisCache) GetHistoryPassword(ctx context.Context, userKey string) ([]string, error) {
	// 系统当前时间
	now := time.Now()
	// 获取当前时间戳（毫秒）
	unixMilli := now.UnixMilli()
	// 获取Zset中的所有member
	result, err := u.cmd.ZRangeByScore(ctx, userKey, &redis.ZRangeBy{
		Min: strconv.FormatInt(unixMilli, 10),
		Max: "+inf",
	}).Result()
	return result, err
}

func (u *UserRedisCache) SetHistoryPassword(ctx context.Context, userKey, historyPassword string) error {
	// 获取当前时间
	currentTime := time.Now()
	// 使用 AddDate 函数计算180天后的时间
	futureTime := currentTime.AddDate(0, 0, 180)
	// 获取时间戳（毫秒）
	futureTimestamp := futureTime.UnixMilli()
	// 执行 Redis 事务操作
	_, err := u.cmd.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		// 插入密码到Redis的ZSet
		err := pipe.ZAdd(ctx, userKey, redis.Z{
			Score:  float64(futureTimestamp),
			Member: historyPassword,
		}).Err()
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (u *UserRedisCache) Set(ctx context.Context, sessionId string) error {
	key := fmt.Sprintf("logout:sessionId:%s", sessionId)
	// 过期时间和长token的过期时间一样
	err := u.cmd.Set(ctx, key, "", time.Hour*24*7).Err()
	if err != nil {
		return fmt.Errorf("session存入redis失败，%w", err)
	}
	return nil
}
