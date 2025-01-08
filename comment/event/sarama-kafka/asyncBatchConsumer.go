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

// AsyncBatchCommentEventConsumer 异步批量消费
type AsyncBatchCommentEventConsumer struct {
	client         sarama.Client // kafka客户端
	commentService service.CommentService
	logger         loggerx.Logger
}

// NewAsyncBatchCommentEventConsumer 创建一个异步批量消费
func NewAsyncBatchCommentEventConsumer(
	client sarama.Client,
	logger loggerx.Logger,
	commentService service.CommentService) *AsyncBatchCommentEventConsumer {
	return &AsyncBatchCommentEventConsumer{
		commentService: commentService,
		client:         client,
		logger:         logger,
	}
}

// StartConsumer 开始消费
func (r *AsyncBatchCommentEventConsumer) StartConsumer() error {
	// 创建消费者组
	cg, err := sarama.NewConsumerGroupFromClient("biz_comment", r.client)
	if err != nil {
		return err
	}
	// 启动一个goroutine，用于异步消费消息
	go func() {
		er := cg.Consume(
			context.Background(),        // 上下文对象
			[]string{"LifeLog_Comment"}, // 消费主题
			saramax.NewBatchHandler[domain.CommentDomain](r.logger, r.Consume),
		)
		if er != nil {
			r.logger.Error("退出消费循环异常", loggerx.Error(er),
				loggerx.String("method:",
					"comment/event/sarama-kafka/AsyncBatchCommentEventConsumer/Start.go"))
		}
	}()
	return nil
}

// Consume 将评论插入数据库
func (r *AsyncBatchCommentEventConsumer) Consume(msg []*sarama.ConsumerMessage, commentDomain []domain.CommentDomain) error {
	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// 确保关闭上下文
	defer cancel()
	// 批量更新阅读计数
	err := r.commentService.BatchCreateComment(ctx, commentDomain)
	if err != nil {
		r.logger.Error("消费评论失败，批量插入数据库失败", loggerx.Error(err))
		return err
	}
	return nil
}
