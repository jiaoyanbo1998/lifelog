package kafkax

import (
	"context"
	"errors"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"sync"
	"time"
)

// Message 消息
type Message struct {
	Data  []byte // 数据
	Topic string // 主题
}

// KafkaAsyncProducer kafka异步生产者
type KafkaAsyncProducer struct {
	writer      *kafka.Writer  // kafka写入器
	Data        chan Message   // 发送数据的管道
	timeout     time.Duration  // 超时时间
	closeChan   chan struct{}  // 发送关闭信号的管道
	wg          sync.WaitGroup // 用于等待所有消息发送完成
	workerCount int            // goroutine 数量
}

// ProducerOption 是一个函数类型，用户自定义配置
type ProducerOption func(*KafkaAsyncProducer)

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) ProducerOption {
	return func(kap *KafkaAsyncProducer) {
		kap.timeout = timeout
	}
}

// WithWorkerCount 设置 goroutine 数量
func WithWorkerCount(workerCount int) ProducerOption {
	return func(kap *KafkaAsyncProducer) {
		kap.workerCount = workerCount
	}
}

// WithChanSize 设置管道大小
func WithChanSize(chanSize int) ProducerOption {
	return func(kap *KafkaAsyncProducer) {
		kap.Data = make(chan Message, chanSize)
	}
}

// NewKafkaAsyncProducer 初始化kafka-go的Producer
func NewKafkaAsyncProducer(addr []string, opts ...ProducerOption) *KafkaAsyncProducer {
	// 创建写入器
	writer := &kafka.Writer{
		Addr:     kafka.TCP(addr...),  // kafka地址
		Balancer: &kafka.LeastBytes{}, // 负载均衡策略
	}
	// 默认配置
	kap := &KafkaAsyncProducer{
		writer:      writer,
		Data:        make(chan Message, 100), // 默认管道大小，100
		timeout:     time.Second * 5,         // 默认超时时间，5秒
		closeChan:   make(chan struct{}),
		workerCount: 1, // 默认 goroutine 数量，1
	}
	// 用户自定义配置
	for _, opt := range opts {
		opt(kap)
	}
	// 启动多个 goroutine，发送消息
	for i := 0; i < kap.workerCount; i++ {
		kap.wg.Add(1)
		go kap.sendMsg()
	}
	return kap
}

// Send 异步发送Message
func (kap *KafkaAsyncProducer) Send(msg Message) {
	// 将Message发送到管道
	kap.Data <- msg
}

// sendMsg 从管道中读取Message，发送到kafka
func (kap *KafkaAsyncProducer) sendMsg() {
	defer kap.wg.Done() // 确保 goroutine结束时通知WaitGroup
	for {
		select {
		// 监听关闭信号
		case <-kap.closeChan:
			// 关闭写入器，退出循环
			err := kap.writer.Close()
			if err != nil {
				zap.L().Error("关闭 Kafka writer 失败", zap.Error(err))
			}
			return
		// 从管道中读取数据
		case data := <-kap.Data:
			// 封装消息
			msg := kafka.Message{
				Topic: data.Topic,
				Value: data.Data,
			}
			var err error
			// 超时时间，控制 WriteMessages() 写入消息的超时时间
			ctx, cancel := context.WithTimeout(context.Background(), kap.timeout)
			defer cancel()
			// 重试 3 次
			const retry = 3
			for i := 0; i < retry; i++ {
				// 发送消息
				err = kap.writer.WriteMessages(ctx, msg)
				// 记录发送的Message
				zap.L().Info("发送消息", zap.String("message", string(msg.Value)), zap.String("topic", msg.Topic))
				// 发送消息成功，退出循环
				if err == nil {
					break
				}
				// 超时错误 DeadlineExceeded，继续循环
				if errors.Is(err, context.DeadlineExceeded) {
					time.Sleep(time.Millisecond * 250)
					continue
				}
				// 其他错误，记录Message
				zap.L().Error("向 Kafka 发送数据失败", zap.Error(err))
			}
			// 如果重试失败，记录错误Message（可以进一步处理，如存储到数据库）
			if err != nil {
				zap.L().Error("发送Message失败，已重试 3 次",
					zap.String("data", string(data.Data)),
					zap.String("topic", data.Topic))
			}
		}
	}
}

// Stop 停止 Kafka 生产者，优雅关闭
func (kap *KafkaAsyncProducer) Stop() {
	close(kap.closeChan) // 发送关闭信号
	kap.wg.Wait()        // 确保在关闭kafka生产者时，所有的消息都发送完成
}
