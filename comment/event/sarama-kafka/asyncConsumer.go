package saramaKafka

import (
	"context"
	"github.com/IBM/sarama"
	"lifelog-grpc/comment/domain"
	"lifelog-grpc/comment/service"
	"lifelog-grpc/pkg/loggerx"
	"lifelog-grpc/pkg/saramax"
	"time"
)

// AsyncCommentEventConsumer 异步消费
type AsyncCommentEventConsumer struct {
	client         sarama.Client // sarama客户端
	commentService service.CommentService
	logger         loggerx.Logger
}

// NewAsyncCommentEventConsumer 创建一个【评论事件】的异步消费者
func NewAsyncCommentEventConsumer(client sarama.Client, l loggerx.Logger,
	commentService service.CommentService) *AsyncCommentEventConsumer {
	return &AsyncCommentEventConsumer{
		commentService: commentService,
		client:         client,
		logger:         l,
	}
}

// StartConsumer 开始消费【评论事件】
func (r *AsyncCommentEventConsumer) StartConsumer() error {
	// 创建消费者组
	// 参数1：消费者组名称
	// 参数2：sarama客户端
	cg, err := sarama.NewConsumerGroupFromClient("biz_comment", r.client)
	if err != nil {
		r.logger.Error("创建消费者组失败", loggerx.Error(err),
			loggerx.String("method:", "comment/event/sarama-kafka/AsyncCommentEventConsumer/Start"))
		return err
	}
	// 启动一个goroutine，异步消费消息
	go func() {
		// 消费消息
		er := cg.Consume(
			context.Background(),                                          // 上下文对象
			[]string{"LifeLog_Comment"},                                   // 消费主题
			saramax.NewHandler[domain.CommentDomain](r.logger, r.Consume), // 消费者处理函数
		)
		if er != nil {
			r.logger.Error("退出消费循环异常", loggerx.Error(er),
				loggerx.String("method:",
					"comment/event/sarama-kafka/AsyncCommentEventConsumer/Start.go"))
		}
	}()
	return err
}

// Consume 处理【评论事件】，将评论插入数据库
func (r *AsyncCommentEventConsumer) Consume(msg *sarama.ConsumerMessage, commentDomain domain.CommentDomain) error {
	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// 确保关闭上下文
	defer cancel()
	// 批量更新阅读计数
	err := r.commentService.CreateComment(ctx, commentDomain)
	if err != nil {
		r.logger.Error("消费评论失败，批量插入数据库失败", loggerx.Error(err))
		return err
	}
	return nil
}
