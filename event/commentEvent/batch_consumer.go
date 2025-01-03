package commentEvent

import (
	"context"
	"github.com/IBM/sarama"
	"time"
	"lifelog-grpc/comment/domain"
	"lifelog-grpc/comment/service"
	"lifelog-grpc/pkg/loggerx"
	"lifelog-grpc/pkg/saramax"
)

// CommentEventBatchConsumer 处理【"评论插入数据库"事件】的消费者
type CommentEventBatchConsumer struct {
	client         sarama.Client // kafka客户端
	commentService service.CommentService
	logger         loggerx.Logger
}

// NewCommentEventBatchConsumer 创建一个【"评论插入数据库"事件】批量消费者
func NewCommentEventBatchConsumer(
	client sarama.Client,
	l loggerx.Logger,
	commentService service.CommentService) *CommentEventBatchConsumer {
	return &CommentEventBatchConsumer{
		commentService: commentService,
		client:         client,
		logger:         l,
	}
}

// Start 开始消费【"评论插入数据库"事件】
func (r *CommentEventBatchConsumer) Start() error {
	// 创建消费者组
	cg, err := sarama.NewConsumerGroupFromClient("comment", r.client)
	if err != nil {
		return err
	}
	// 启动一个goroutine，用于消费消息
	go func() {
		er := cg.Consume(context.Background(), []string{"insert_comment"},
			saramax.NewBatchHandler[CommentEvent](r.logger, r.Consume))
		if er != nil {
			r.logger.Error("退出了消费循环异常", loggerx.Error(err))
		}
	}()
	return err
}

// Consume 处理接收到的【文章阅读事件】，更新文章的阅读计数
func (r *CommentEventBatchConsumer) Consume(
	msg []*sarama.ConsumerMessage, ts []CommentEvent) error {
	var comments []domain.CommentDomain
	// 遍历消息，提取文章ID和用户ID
	for _, t := range ts {
		comments = append(comments, domain.CommentDomain{
			Biz:      "comment",
			UserId:   t.UserId,
			BizId:    t.BizId,
			Content:  t.Content,
			ParentId: t.ParentId.Int64,
			RootId:   t.RootId.Int64,
		})
	}
	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// 确保关闭上下文
	defer cancel()
	// ”批量插入“评论
	err := r.commentService.BatchCreateComment(ctx, comments)
	if err != nil {
		r.logger.Error("批量插入评论失败", loggerx.Error(err))
	}
	return nil
}
