package lifeLogEvent

import (
	"context"
	"github.com/IBM/sarama"
	"time"
	"lifelog-grpc/interactive/service"
	"lifelog-grpc/pkg/loggerx"
	"lifelog-grpc/pkg/saramax"
)

// AsyncReadEventConsumer 【文章阅读事件】的消费者
type AsyncReadEventConsumer struct {
	client             sarama.Client
	interactiveService service.InteractiveService
	logger             loggerx.Logger
}

// NewAsyncReadEventConsumer 创建一个异步的【文章阅读事件】消费者
func NewAsyncReadEventConsumer(client sarama.Client, l loggerx.Logger,
	interactiveService service.InteractiveService) *AsyncReadEventConsumer {
	return &AsyncReadEventConsumer{
		interactiveService: interactiveService,
		client:             client,
		logger:             l,
	}
}

// Start 开始异步消费文章【文章阅读事件】
func (r *AsyncReadEventConsumer) Start() error {
	// 创建消费者组
	cg, err := sarama.NewConsumerGroupFromClient("interactive", r.client)
	if err != nil {
		return err
	}
	// 启动一个goroutine，用于异步消费消息
	go func() {
		er := cg.Consume(
			context.Background(),
			[]string{"read_lifeLog"},
			saramax.NewHandler[ReadEvent](r.logger, r.Consume),
		)
		if er != nil {
			r.logger.Error("退出了消费循环异常", loggerx.Error(err))
		}
	}()
	return nil
}

// Consume 处理接收到的【文章阅读事件】，更新文章的阅读计数
func (r *AsyncReadEventConsumer) Consume(msg *sarama.ConsumerMessage, evt ReadEvent) error {
	go func() {
		// 创建带超时的上下文
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		// 异步更新阅读计数
		err := r.interactiveService.IncreaseReadCount(ctx,
			"lifeLog", evt.LifeLogId, evt.UserId)
		if err != nil {
			r.logger.Error("更新阅读计数失败", loggerx.Error(err))
		}
	}()
	return nil
}
