package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type RankingCache interface {
	GetTopN(ctx context.Context, number int64) ([]string, error)
}

type RankingCacheRedis struct {
	cmd redis.Cmdable
}

func NewRankingCacheRedis(cmd redis.Cmdable) RankingCache {
	return &RankingCacheRedis{
		cmd: cmd,
	}
}

func (r *RankingCacheRedis) GetTopN(ctx context.Context, number int64) (
	[]string, error) {
	// 获取前n篇文章的title
	res, err := r.cmd.ZRevRange(ctx,
		"ranking_article", 0, number).Result()
	if err != nil {
		return nil, err
	}
	return res, nil
}
