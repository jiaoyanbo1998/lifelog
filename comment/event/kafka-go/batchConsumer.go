package ioc

import (
	"context"
	"github.com/spf13/viper"
	"lifelog-grpc/comment/domain"
	"lifelog-grpc/comment/service"
	"lifelog-grpc/pkg/kafkax"
	"lifelog-grpc/pkg/loggerx"
	"time"
)

type CommentBatchConsumer struct {
	commentService service.CommentService
	logger         loggerx.Logger
}

func NewCommentBatchConsumer(commentService service.CommentService,
	logger loggerx.Logger) *CommentBatchConsumer {
	return &CommentBatchConsumer{
		commentService: commentService,
		logger:         logger,
	}
}

func (c *CommentBatchConsumer) Start(logger loggerx.Logger) {
	type kafkaConfig struct {
		Addr    []string      `yaml:"addr"`
		GroupId string        `yaml:"groupId"`
		Topic   string        `yaml:"topic"`
		Timeout time.Duration `yaml:"timeout"`
		Batch   int           `yaml:"batch"`
	}
	var kCfg kafkaConfig
	err := viper.UnmarshalKey("commentKafka", &kCfg)
	if err != nil {
		logger.Error("加载配置文件失败", loggerx.Error(err))
		panic(err)
	}
	// 创建批量消费者
	consumer := kafkax.NewKafkaAsyncBatchConsumer[domain.CommentDomain](kCfg.Addr, kCfg.GroupId, kCfg.Topic,
		kCfg.Timeout*time.Second, kCfg.Batch, c.handler())
	// 关闭消费者
	defer consumer.Stop()
	// 异步消费
	go func() {
		consumer.ReadAndProcessMsg()
	}()
}

// 定义 handler 函数
func (c *CommentBatchConsumer) handler() func(valS []domain.CommentDomain) error {
	return func(valS []domain.CommentDomain) error {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		err := c.commentService.BatchCreateComment(ctx, valS)
		if err != nil {
			c.logger.Error("批量消费失败", loggerx.Error(err),
				loggerx.String("method:", "CommentConsumer:Start:handler"))
			return err
		}
		return nil
	}
}
