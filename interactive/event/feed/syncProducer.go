package feed

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

// FeedEvent 发送评论Feed流
type FeedEvent struct {
	UserId     int64  `json:"user_id"`
	Content    string `json:"content"`     // 事件的内容(业务内容)
	Type       string `json:"type"`        // 事件的类型，根据不同的事件类型，调用不同的handler
	CreateTime int64  `json:"create_time"` // 创建时间
}

// ProduceInteractiveEventFeed 生产InteractiveEvent的feed流
func (s *SyncProducer) ProduceInteractiveEventFeed(feedEvent FeedEvent) error {
	val, err := json.Marshal(feedEvent)
	// json序列化失败
	if err != nil {
		s.logger.Error("json序列化失败", loggerx.Error(err),
			loggerx.String("method:", "interactive:event:feed:sarama-kafka:SyncProducer:ProduceCommentEvent"))
		return err
	}
	_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
		// 设置主题的名字
		Topic: "LifeLog_feed_event",
		// 设置消息内容，消息的内容必须为[]byte类型
		Value: sarama.ByteEncoder(val),
	})
	if err != nil {
		s.logger.Error("消息发送失败", loggerx.Error(err),
			loggerx.String("method:", "interactive:event:feed:sarama-kafka:SyncProducer:ProduceCommentEvent"))
		return err
	}
	return nil
}
