package event

import (
	"context"
	"github.com/IBM/sarama"
	"lifelog-grpc/lifeLog/service"
	"lifelog-grpc/pkg/loggerx"
	"lifelog-grpc/pkg/saramax"
	"time"
)

// AsyncLifeLogEventConsumer 异步消费
type AsyncLifeLogEventConsumer struct {
	client         sarama.Client // sarama客户端
	logger         loggerx.Logger
	lifeLogService service.LifeLogService
}

// NewAsyncLifeLogEventConsumer 创建一个【lifeLog缓存预加载事件】的异步消费者
func NewAsyncLifeLogEventConsumer(client sarama.Client, l loggerx.Logger,
	lifeLogService service.LifeLogService) *AsyncLifeLogEventConsumer {
	return &AsyncLifeLogEventConsumer{
		client:         client,
		logger:         l,
		lifeLogService: lifeLogService,
	}
}

// StartConsumer 开始消费【lifeLog缓存预加载事件】
func (r *AsyncLifeLogEventConsumer) StartConsumer() error {
	// 创建消费者组
	// 参数1：消费者组名称
	// 参数2：sarama客户端
	cg, err := sarama.NewConsumerGroupFromClient("lifeLog_List", r.client)
	if err != nil {
		r.logger.Error("创建消费者组失败", loggerx.Error(err),
			loggerx.String("method:", "lifeLog:event:AsyncConsumer:StartConsumer"))
		return err
	}
	// 启动一个goroutine，异步消费消息
	go func() {
		// 消费消息
		er := cg.Consume(
			context.Background(),                                  // 上下文对象
			[]string{"lifeLog_List"},                              // 消费主题
			saramax.NewHandler[EventLifeLog](r.logger, r.Consume), // 消费者处理函数
		)
		if er != nil {
			r.logger.Error("退出消费循环异常", loggerx.Error(er),
				loggerx.String("method:", "lifeLog:event:AsyncConsumer:StartConsumer"))
		}
	}()
	return err
}

// Consume 处理【lifeLog缓存预加载事件】
func (r *AsyncLifeLogEventConsumer) Consume(msg *sarama.ConsumerMessage, lifeLogEvent EventLifeLog) error {
	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// 确保关闭上下文
	defer cancel()
	// lifeLog缓存预加载
	_, err := r.lifeLogService.SearchByAuthorId(ctx, lifeLogEvent.AuthorId, lifeLogEvent.Limit, lifeLogEvent.Offset)
	if err != nil {
		r.logger.Error("消费lifeLog缓存预加载事件失败，插入redis失败", loggerx.Error(err))
		return err
	}
	return nil
}
