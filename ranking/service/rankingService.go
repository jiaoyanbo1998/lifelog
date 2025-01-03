package service

import (
	"context"
	"lifelog-grpc/ranking/repository"
)

type RankingService interface {
	TopN(ctx context.Context, number int64) ([]string, error)
}

type RankingServiceV1 struct {
	rankingRepository repository.RankingRepository
}

func NewRankingService(rankingRepository repository.RankingRepository) RankingService {
	return &RankingServiceV1{
		rankingRepository: rankingRepository,
	}
}

// TopN 获取热度最高的前n篇文章
func (r *RankingServiceV1) TopN(ctx context.Context, number int64) (
	[]string, error) {
	return r.rankingRepository.TopN(ctx, number)
}
