package ioc

import (
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"lifelog-grpc/job"
	"lifelog-grpc/pkg/loggerx"
	"lifelog-grpc/ranking/service"
	"time"
)

// InitRankingJob 初始化热榜任务
func InitRankingJob(rankingService service.RankingService, logger loggerx.Logger,
	cmd redis.Cmdable) job.Job {
	// ranking，热榜任务的执行最长时间
	//   超过30s，上下文就会被取消，热榜任务也就会被中断
	t := time.Second * 30
	rankingJob := job.NewRankingJob(rankingService, logger, t)
	return rankingJob
}

// InitCronRankingJob 初始化热榜定时任务
func InitCronRankingJob(logger loggerx.Logger, rankingJob job.Job) *cron.Cron {
	// 创建一个cron，允许秒级调度
	cron := cron.New(cron.WithSeconds())
	// 创建热榜任务适配器，将我们编写大的热榜任务，转化为cron的定时任务
	adapter := job.NewCronRankingJobAdapter(logger, rankingJob)
	// 每3分钟调度一次
	_, err := cron.AddJob("@every 3m", adapter)
	// 调度失败
	if err != nil {
		panic(err)
	}
	return cron
}
