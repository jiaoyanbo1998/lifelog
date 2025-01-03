package repository

import (
	"context"
	"lifelog-grpc/ranking/repository/cache"
)

type RankingRepository interface {
	TopN(ctx context.Context, number int64) ([]string, error)
}

type RankingRepositoryV1 struct {
	rankingCache cache.RankingCache
}

func NewRankingRepository(rankingCache cache.RankingCache) RankingRepository {
	return &RankingRepositoryV1{
		rankingCache: rankingCache,
	}
}

func (r *RankingRepositoryV1) TopN(ctx context.Context, number int64) (
	[]string, error) {
	return r.rankingCache.GetTopN(ctx, number)
}
