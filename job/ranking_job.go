package job

import (
	"context"
	"time"
	"lifelog-grpc/pkg/loggerx"
	"lifelog-grpc/ranking/service"
)

// RankingJob 热榜
type RankingJob struct {
	rankingService service.RankingService
	timeout        time.Duration
	number         int64
}

func NewRankingJob(rankingService service.RankingService,
	logger loggerx.Logger, timeout time.Duration) *RankingJob {
	return &RankingJob{
		rankingService,
		timeout,
		10,
	}
}

func (r *RankingJob) Name() string {
	return "ranking"
}

func (r *RankingJob) Run() (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	// 获取热榜（热度最高的前n篇文章）
	return r.rankingService.TopN(ctx, r.number)
}
