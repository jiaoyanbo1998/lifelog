package kafkax

import (
	"context"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"time"
)

// KafkaBatchConsumer kafka批量消费
type KafkaBatchConsumer struct {
	reader    *kafka.Reader // kafka读取器
	batchSize int           // 每一批的大小
	timeout   time.Duration // 消费超时时间
	closeChan chan struct{} // 关闭信号通道
}

// NewKafkaBatchConsumer 创建一个Kafka消费者
func NewKafkaBatchConsumer(brokers []string, groupId, topic string,
	batchSize int, timeout time.Duration) *KafkaBatchConsumer {
	// 创建kafka读取器
	reader := kafka.NewReader(kafka.ReaderConfig{
		// Brokers是一个切片，每个元素都是一个节点的地址(kafka地址)
		Brokers:  brokers,
		GroupID:  groupId, // 消费组Id
		Topic:    topic,   // 消费的主题
		MaxBytes: 10e6,    // 每次读取的最大字节数
	})
	// 创建消费者实例
	kbc := &KafkaBatchConsumer{
		reader:    reader,
		batchSize: batchSize,
		timeout:   timeout,
		closeChan: make(chan struct{}),
	}
	// 启动一个goroutine，读取数据
	go kbc.readMsg()
	return kbc
}

// readMsg 从kafka读取消息
func (kbc *KafkaBatchConsumer) readMsg() {
	// 消费循环
	for {
		// 每次for，都创建一个新的超时上下文
		//ctx, cancel := context.WithTimeout(context.Background(), kbc.timeout)
		select {
		// 监听关闭信号
		case <-kbc.closeChan:
			zap.L().Info("批量消费者关闭")
			kbc.reader.Close()
			// 退出循环
			return
		default:
			// 存储本次读取到的消息
			var messages []kafka.Message
			// 批量读取消息
			for i := 0; i < kbc.batchSize; i++ {
				// 从kafka中读取消息
				m, err := kbc.reader.FetchMessage(context.Background())
				if err != nil {
					zap.L().Error("从kafka中读取数据失败", zap.Error(err))
					// 跳过本次循环
					continue
				}
				messages = append(messages, m)
			}
			// 如果没有消息，继续下一次循环
			if len(messages) == 0 {
				continue
			}
			// 批量提交，已经消费的消息
			for _, message := range messages {
				// 标记消息已消费
				err := kbc.reader.CommitMessages(context.Background(), message)
				if err != nil {
					zap.L().Error("标记消息失败", zap.Error(err))
				} else {
					zap.L().Info("消息已消费",
						zap.String("topic", message.Topic),
						zap.Int64("offset", message.Offset))
				}
			}
		}
		// 释放超时上下文
		//cancel()
	}
}

// Stop 停止Kafka消费者，优雅关闭
func (kbc *KafkaBatchConsumer) Stop() {
	close(kbc.closeChan) // 发送关闭信号，停止消费者
}
