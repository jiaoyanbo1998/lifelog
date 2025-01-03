package lifeLogEvent

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"lifelog-grpc/pkg/loggerx"
)

// SaramaAsyncProducer 异步生产者
type SaramaAsyncProducer struct {
	producer sarama.AsyncProducer
	logger   loggerx.Logger
}

// NewSaramaAsyncProducer 使用构造函数，创建一个异步生产者
func NewSaramaAsyncProducer(producer sarama.AsyncProducer, l loggerx.Logger) Producer {
	return &SaramaAsyncProducer{
		producer: producer,
		logger:   l,
	}
}

// ProduceReadEvent 使用异步生产者，将消息发送到kafka的主题
func (s *SaramaAsyncProducer) ProduceReadEvent(evt ReadEvent) error {
	// json序列化，将数据转为json格式，json格式表现形式是[]byte
	val, err := json.Marshal(evt)
	// json序列化失败
	if err != nil {
		s.logger.Error("json序列化失败", loggerx.Error(err))
		return err
	}
	// 获取生产者的输入通道
	messagesChan := s.producer.Input()
	// 将消息发送到kafka的主题
	messagesChan <- &sarama.ProducerMessage{
		// 主题的名字【重要】
		Topic: "read_lifeLog",
		// 消息的key（会通过计算key的哈希值，将相同的key的消息发送到同一个分区）
		Key: sarama.StringEncoder("biz-lifeLog"),
		// 消息内容【重要】
		Value: sarama.ByteEncoder(val),
		// 消息头，在生产者和消费者之间传递（可选）
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("key"),
				Value: []byte("value"),
			},
		},
		// 额外的信息（可选）
		Metadata: map[string]any{"metadataKey": "metadataValue"},
	}
	// 处理结果（监听发送成功了，还是发送失败了）
	go func() {
		// 哪一个case成功，就执行哪一个case
		// 所有的case都成功，随机执行一个case
		// 所有的case都失败，就执行default语句，
		//		没有default语句，就阻塞等待某一个case成功
		select {
		// 发送失败
		case er := <-s.producer.Errors():
			s.logger.Error("消息发送失败", loggerx.Error(er))
		// 发送成功
		case suc := <-s.producer.Successes():
			// json序列化
			bytes, er := json.Marshal(suc.Value)
			if er != nil {
				s.logger.Error("json序列化失败", loggerx.Error(err))
			}
			s.logger.Info("发送成功",
				loggerx.String("value：", string(bytes)))
		}
	}()
	return nil
}
