package saramaKafka

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"lifelog-grpc/comment/domain"
	"lifelog-grpc/pkg/loggerx"
)

// AsyncProducer 异步生产者
type AsyncProducer struct {
	producer sarama.AsyncProducer
	logger   loggerx.Logger
}

// NewAsyncProducer 创建一个异步生产者
func NewAsyncProducer(producer sarama.AsyncProducer, logger loggerx.Logger) *AsyncProducer {
	return &AsyncProducer{
		producer: producer,
		logger:   logger,
	}
}

// Close 关闭生产者
func (as *AsyncProducer) Close() error {
	return as.producer.Close()
}

// ProduceCommentEvent 生产commentEvent
func (as *AsyncProducer) ProduceCommentEvent(commentDomain domain.CommentDomain) error {
	// json序列化（将数据转为json格式，json格式表现形式是[]byte）
	val, err := json.Marshal(commentDomain)
	// json序列化失败
	if err != nil {
		as.logger.Error("json序列化失败", loggerx.Error(err),
			loggerx.String("method:", "comment:event:sarama-kafka:AsyncProducer:ProduceCommentEvent"))
		return err
	}
	// 获取生产者的输入通道
	messagesChan := as.producer.Input()
	// 将消息发送到kafka的主题
	messagesChan <- &sarama.ProducerMessage{
		// 主题的名字
		Topic: "LifeLog_comment_event",
		// 消息的key（通过计算key的哈希值，将相同的key的消息发送到同一个分区）
		Key: sarama.StringEncoder("biz-LifeLog"),
		// 消息内容
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
		select {
		// 发送失败
		case er := <-as.producer.Errors():
			as.logger.Error("消息发送失败", loggerx.Error(er))
		// 发送成功
		case suc := <-as.producer.Successes():
			// json序列化
			bytes, er := json.Marshal(suc.Value)
			if er != nil {
				as.logger.Error("json序列化失败", loggerx.Error(err))
				return
			}
			as.logger.Info("发送成功",
				loggerx.String("value：", string(bytes)))
		}
	}()
	return nil
}
