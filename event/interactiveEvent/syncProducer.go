package interactiveEvent

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"lifelog-grpc/pkg/loggerx"
)

// SyncProducer 同步生产者
type SyncProducer struct {
	producer sarama.SyncProducer
	logger   loggerx.Logger
}

// NewSyncProducer 创建一个同步生产者
func NewSyncProducer(producer sarama.SyncProducer, logger loggerx.Logger) *SyncProducer {
	return &SyncProducer{
		producer: producer,
		logger:   logger,
	}
}

// Close 关闭生产者
func (s *SyncProducer) Close() error {
	return s.producer.Close()
}

// ProduceInteractiveEvent 生产InteractiveEvent
func (s *SyncProducer) ProduceInteractiveEvent(interactiveEvent EventInteractive) error {
	// 将InteractiveEvent序列化为json（[]byte格式的json对象）
	val, err := json.Marshal(interactiveEvent)
	// json序列化失败
	if err != nil {
		s.logger.Error("json序列化失败", loggerx.Error(err),
			loggerx.String("method:", "event:interactiveEvent:SyncProducer:ProduceLifeLogEvent"))
		return err
	}
	// 将InteractiveEvent发送到Kafka的topic
	// 返回值：分区，offset(消息在分区中的位置)，错误信息
	_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
		// 设置主题的名字
		Topic: "lifeLog_interactive_event_liked_count",
		// 设置消息内容，消息的内容必须为[]byte类型
		Value: sarama.ByteEncoder(val),
	})
	if err != nil {
		s.logger.Error("消息发送失败", loggerx.Error(err),
			loggerx.String("method:", "event:interactiveEvent:SyncProducer:ProduceLifeLogEvent"))
		return err
	}
	return nil
}
