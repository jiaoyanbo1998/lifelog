package commentEvent

import (
	"encoding/json"
	"github.com/IBM/sarama"
)

// SaramaSyncProducer 同步生产者
type SaramaSyncProducer struct {
	producer sarama.SyncProducer
}

// NewSaramaSyncProducer 创建一个同步生产者
func NewSaramaSyncProducer(producer sarama.SyncProducer) Producer {
	return &SaramaSyncProducer{
		producer: producer,
	}
}

// ProduceCommentEvent 使用同步生产者，将消息发送到Kafka的主题
func (s *SaramaSyncProducer) ProduceCommentEvent(evt CommentEvent) error {
	// json反序列化，文章阅读事件
	val, err := json.Marshal(evt)
	// json反序列化失败
	if err != nil {
		return err
	}
	// 将消息发送到Kafka的topic
	// 		返回值：分区，offset（消息在分区中的位置），错误信息
	_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
		// 设置主题的名字 【重要】
		Topic: "insert_comment",
		// 设置消息内容，val必须为[]byte类型【重要】
		Value: sarama.ByteEncoder(val),
	})
	return err
}
