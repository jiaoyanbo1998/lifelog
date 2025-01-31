package likedEvent

import (
	"context"
	"github.com/IBM/sarama"
	"lifelog-grpc/interactive/service"
	"lifelog-grpc/pkg/loggerx"
	"lifelog-grpc/pkg/saramax"
	"time"
)

// AsyncLikedEventConsumer 异步消费
type AsyncLikedEventConsumer struct {
	client             sarama.Client // sarama客户端
	interactiveService service.InteractiveService
	logger             loggerx.Logger
}

// NewAsyncLikedEventConsumer 创建一个【点赞事件】的异步消费者
func NewAsyncLikedEventConsumer(client sarama.Client, l loggerx.Logger,
	interactiveService service.InteractiveService) *AsyncLikedEventConsumer {
	return &AsyncLikedEventConsumer{
		client:             client,
		logger:             l,
		interactiveService: interactiveService,
	}
}

// StartConsumer 开始消费【点赞事件】
func (r *AsyncLikedEventConsumer) StartConsumer() error {
	// 创建消费者组
	// 参数1：消费者组名称
	// 参数2：sarama客户端
	cg, err := sarama.NewConsumerGroupFromClient("lifeLog_liked", r.client)
	if err != nil {
		r.logger.Error("创建消费者组失败", loggerx.Error(err),
			loggerx.String("method:", "interactive/event/likedEvent/asyncConsumer/StartConsumer"))
		return err
	}
	// 启动一个goroutine，异步消费消息
	go func() {
		// 消费消息
		er := cg.Consume(
			context.Background(),                                      // 上下文对象
			[]string{"lifeLog_interactive_event_liked_count"},         // 消费主题
			saramax.NewHandler[EventInteractive](r.logger, r.Consume), // 消费者处理函数
		)
		if er != nil {
			r.logger.Error("退出消费循环异常", loggerx.Error(er),
				loggerx.String("method:",
					"interactive/event/likedEvent/asyncConsumer/StartConsumer"))
		}
	}()
	return err
}

// Consume 处理【评论事件】，将评论插入数据库
func (r *AsyncLikedEventConsumer) Consume(msg *sarama.ConsumerMessage, interactiveEvent EventInteractive) error {
	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// 确保关闭上下文
	defer cancel()
	// 更新点赞计数
	err := r.interactiveService.IncreaseReadCount(ctx, interactiveEvent.Biz, interactiveEvent.BizId, interactiveEvent.UserId)
	if err != nil {
		r.logger.Error("消费点赞时间失败，插入数据库失败", loggerx.Error(err))
		return err
	}
	return nil
}
