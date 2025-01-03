package lifeLogEvent

import (
	"context"
	"github.com/IBM/sarama"
	"time"
	"lifelog-grpc/interactive/service"
	"lifelog-grpc/pkg/loggerx"
	"lifelog-grpc/pkg/saramax"
)

// ReadEventBatchConsumer 处理【文章阅读事件】的消费者
type ReadEventBatchConsumer struct {
	client             sarama.Client // kafka客户端
	interactiveService service.InteractiveService
	logger             loggerx.Logger
}

// NewReadEventBatchConsumer 创建一个【文章阅读事件】批量消费者
func NewReadEventBatchConsumer(
	client sarama.Client,
	l loggerx.Logger,
	interactiveService service.InteractiveService) *ReadEventBatchConsumer {
	return &ReadEventBatchConsumer{
		interactiveService: interactiveService,
		client:             client,
		logger:             l,
	}
}

// Start 开始消费【文章阅读事件】
func (r *ReadEventBatchConsumer) Start() error {
	// 创建消费者组
	cg, err := sarama.NewConsumerGroupFromClient("interactive", r.client)
	if err != nil {
		return err
	}
	// 启动一个goroutine，用于异步消费消息
	go func() {
		er := cg.Consume(context.Background(), []string{"read_lifeLog"},
			saramax.NewBatchHandler[ReadEvent](r.logger, r.Consume))
		if er != nil {
			r.logger.Error("退出了消费循环异常", loggerx.Error(er),
				loggerx.String("method:", "event/lifeLogEvent/batch_consumer.go"))
		}
	}()
	return nil
}

// Consume 处理接收到的【文章阅读事件】，更新文章的阅读计数
func (r *ReadEventBatchConsumer) Consume(msg []*sarama.ConsumerMessage, ts []ReadEvent) error {
	LifeLogIds := make([]int64, 0, len(ts))
	UserIds := make([]int64, 0, len(ts))
	// 遍历消息，提取文章ID和用户ID
	for _, t := range ts {
		LifeLogIds = append(LifeLogIds, t.LifeLogId)
		UserIds = append(UserIds, t.UserId)
	}
	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// 确保关闭上下文
	defer cancel()
	// 批量更新阅读计数
	err := r.interactiveService.BatchInteractiveReadCount(ctx,
		"lifeLog", LifeLogIds, UserIds)
	if err != nil {
		r.logger.Error("批量增加阅读计数失败", loggerx.Error(err))
	}
	return nil
}
