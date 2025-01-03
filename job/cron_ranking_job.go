package job

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
	"lifelog-grpc/pkg/loggerx"
)

// CronRankingJobAdapter 定时任务
type CronRankingJobAdapter struct {
	job        Job // 自定义的job
	logger     loggerx.Logger
	summaryVec *prometheus.SummaryVec
}

func NewCronRankingJobAdapter(
	l loggerx.Logger, j Job) *CronRankingJobAdapter {
	// 初始化监控指标
	vec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "jyb",
		Subsystem: "webook",
		Name:      "ranking_job",
		Help:      "监控定时任务的执行时间，是否执行成功",
		Objectives: map[float64]float64{
			0.5:  0.05,  // 50%的任务在0.05s内完成
			0.7:  0.02,  // 70%的任务在0.02s内完成
			0.9:  0.01,  // 90%的任务在0.01s内完成
			0.95: 0.005, // 95%的任务在0.005s内完成
		},
		// 监控指标
	}, []string{"job_name", "success"})
	return &CronRankingJobAdapter{
		logger:     l,
		job:        j,
		summaryVec: vec,
	}
}

// Run 适配器模式（实现Cron的Job接口）
//		将自定义的Job转为Cron的Job（定时任务）
func (c *CronRankingJobAdapter) Run() {
	// 任务开始时间
	start := time.Now()
	// 定时任务是否执行成功
	var success = "true"
	defer func() {
		// 任务结束时间
		end := time.Now()
		// 更新监控指标的值
		c.summaryVec.WithLabelValues(c.job.Name(), success).
			// 记录执行时间，end - start
			Observe(float64(end.Sub(start)))
	}()
	// 执行任务
	_, err := c.job.Run()
	// 记录错误日志
	if err != nil {
		c.logger.Error("定时任务执行失败",
			loggerx.String("name:", c.job.Name()),
			loggerx.String("method:", "CronRankingJob:Run"),
			loggerx.Error(err))
		success = "false"
	}
}
