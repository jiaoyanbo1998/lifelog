package saramax

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"lifelog-grpc/pkg/loggerx"
	"time"
)

// Handler 泛型结构体，用于处理kafka消费者信息
type Handler[T any] struct {
	logger loggerx.Logger
	// 消费处理函数
	MessageHandlerFunc func(msg *sarama.ConsumerMessage, t T) error
}

// NewHandler 使用构造方法，创建一个Handler
func NewHandler[T any](l loggerx.Logger,
	fn func(msg *sarama.ConsumerMessage, t T) error) *Handler[T] {
	return &Handler[T]{
		logger:             l,
		MessageHandlerFunc: fn,
	}
}

// ConsumeClaim 消费Kafka消息
// 参数：session 消费者组会话
// 参数：claim 消费者组的分区
func (h *Handler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// 获取分区中的所有消息
	messages := claim.Messages()
	// 循环处理消息
	for message := range messages {
		// 消息中的value
		var val T
		// json反序列化消息
		err := json.Unmarshal(message.Value, &val)
		if err != nil {
			// json反序列化失败，记录日志
			h.logger.Error("反序列化消息体失败",
				loggerx.String("topic", message.Topic),
				loggerx.Int("partition", int(message.Partition)),
				loggerx.Int("offset", int(message.Offset)),
				loggerx.Error(err))
			// 不中断，继续下一个
			// 反序列化失败，不需要重试，因为数据格式本身就是错的，重试多少次还是会出错
			continue
		}
		// 处理消息（带有重试）
		retryCount := 3
		for i := 0; i < retryCount; i++ {
			// 消息处理函数
			err = h.MessageHandlerFunc(message, val)
			if err == nil {
				// 处理成功，标记消息为已消费
				session.MarkMessage(message, "")
				break
			}
			// 处理失败，记录日志并重试
			if i == retryCount-1 {
				// 重试次数达到上限，记录错误日志
				h.logger.Error("重试次数达到上限",
					loggerx.String("topic", message.Topic),
					loggerx.Int("partition", int(message.Partition)),
					loggerx.Int("offset", int(message.Offset)),
					loggerx.Error(err))
				// 可以将消息发送到死信队列以便后续处理
			} else {
				// 记录重试日志
				h.logger.Warn("处理消息失败，准备重试",
					loggerx.String("topic", message.Topic),
					loggerx.Int("partition", int(message.Partition)),
					loggerx.Int("offset", int(message.Offset)),
					loggerx.Error(err))
				// 指数退避
				// 休眠多长时间：1s、2s、4s
				// time.Duration(i+1) == i+1
				time.Sleep(time.Duration(i+1) * time.Second)
			}
		}
	}
	return nil
}

// Setup 消费者组会话开始前进行初始化工作
func (h *Handler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup 消费者组会话结束后进行清理工作
func (h *Handler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}
