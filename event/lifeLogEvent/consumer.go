package lifeLogEvent

import (
	"context"
	"github.com/IBM/sarama"
	"lifelog-grpc/interactive/service"
	"lifelog-grpc/pkg/loggerx"
	"lifelog-grpc/pkg/saramax"
	"time"
)

// ReadEventConsumer 【文章阅读事件】的消费者
type ReadEventConsumer struct {
	client             sarama.Client
	interactiveService service.InteractiveService
	logger             loggerx.Logger
}

func NewReadEventConsumer(client sarama.Client, l loggerx.Logger,
	interactiveService service.InteractiveService) *ReadEventConsumer {
	return &ReadEventConsumer{
		interactiveService: interactiveService,
		client:             client,
		logger:             l,
	}
}

// Start 开始消费文章【文章阅读事件】
func (r *ReadEventConsumer) Start() error {
	// 创建消费者组
	cg, err := sarama.NewConsumerGroupFromClient("interactive", r.client)
	if err != nil {
		return err
	}
	// 启动一个goroutine，用于消费消息
	go func() {
		er := cg.Consume(
			context.Background(),
			[]string{"read_LifeLog"},
			saramax.NewHandler[ReadEvent](r.logger, r.Consume),
		)
		if er != nil {
			r.logger.Error("退出了消费循环异常", loggerx.Error(er),
				loggerx.String("method:", "event/LifeLogEvent/consumer.go"))
		}
	}()
	return nil
}

// Consume 处理接收到的【文章阅读事件】，更新文章的阅读计数
func (r *ReadEventConsumer) Consume(msg *sarama.ConsumerMessage, evt ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := r.interactiveService.IncreaseReadCount(ctx,
		"lifeLog", evt.LifeLogId, evt.UserId)
	return err
}
