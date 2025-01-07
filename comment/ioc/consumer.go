package ioc

import (
	"context"
	"github.com/spf13/viper"
	commentv1 "lifelog-grpc/api/proto/gen/comment/v1"
	"lifelog-grpc/comment/grpc"
	"lifelog-grpc/pkg/kafkax"
	"lifelog-grpc/pkg/loggerx"
	"time"
)

type CommentConsumer struct {
	commentServiceGRPCService *grpc.CommentServiceGRPCService
	logger                    loggerx.Logger
}

func NewCommentConsumer(commentServiceGRPCService *grpc.CommentServiceGRPCService, logger loggerx.Logger) *CommentConsumer {
	return &CommentConsumer{
		commentServiceGRPCService: commentServiceGRPCService,
		logger:                    logger,
	}
}

type comment struct {
	UserId   int64  `json:"user_id"`
	Biz      string `json:"biz"`
	BizId    int64  `json:"biz_id"`
	Content  string `json:"content"`
	ParentId int64  `json:"parent_id"`
	RootId   int64  `json:"root_id"`
}

func (c *CommentConsumer) Start(logger loggerx.Logger) {
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
	consumer := kafkax.NewKafkaAsyncBatchConsumer[comment](kCfg.Addr, kCfg.GroupId, kCfg.Topic,
		kCfg.Timeout*time.Second, kCfg.Batch, c.handler())
	// 关闭消费者
	defer consumer.Stop()
	// 异步消费
	go func() {
		consumer.ReadAndProcessMsg()
	}()
}

// 定义 handler 函数
func (c *CommentConsumer) handler() func(valS []comment) error {
	return func(valS []comment) error {
		// 将[]comment，转为，[]*Comment
		cts := make([]*commentv1.Comment, 0, len(valS))
		for _, val := range valS {
			cts = append(cts, &commentv1.Comment{
				UserId:   val.UserId,
				Biz:      val.Biz,
				BizId:    val.BizId,
				Content:  val.Content,
				ParentId: val.ParentId,
				RootId:   val.RootId,
			})
		}
		_, err := c.commentServiceGRPCService.BatchCreateComment(context.Background(),
			&commentv1.BatchCreateCommentRequest{
				Comment: cts,
			})
		if err != nil {
			c.logger.Error("批量消费失败", loggerx.Error(err),
				loggerx.String("method:", "CommentConsumer:Start:handler"))
			return err
		}
		return nil
	}
}
