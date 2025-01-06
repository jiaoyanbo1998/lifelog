package kafkax

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// KafkaSyncConsumer 异步消费
type KafkaSyncConsumer struct {
	reader  *kafka.Reader // kafka读取器
	timeout time.Duration // 超时时间
}

// Option 是一个函数类型，用于设置 KafkaSyncConsumer 的选项
type Option func(*kafka.ReaderConfig)

// NewKafkaConsumer 创建一个kafka消费者
func NewKafkaConsumer(brokers []string, groupId, topic string, timeout time.Duration, opts ...Option) *KafkaSyncConsumer {
	// 默认配置
	config := kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupId,
		Topic:    topic,
		MinBytes: 10e3, // 默认 最小字节数 10KB
		MaxBytes: 10e6, // 默认 最大字节数 10MB
	}
	// 用户自定义配置
	for _, opt := range opts {
		opt(&config)
	}
	// 创建 KafkaSyncConsumer
	ksc := &KafkaSyncConsumer{
		reader:  kafka.NewReader(config),
		timeout: timeout,
	}
	// 启动一个goroutine，读取数据
	go ksc.readMsg()
	return ksc
}

// WithMinBytes 设置每次读取的最小字节数
func WithMinBytes(minBytes int) Option {
	return func(cfg *kafka.ReaderConfig) {
		cfg.MinBytes = minBytes
	}
}

// WithMaxBytes 设置每次读取的最大字节数
func WithMaxBytes(maxBytes int) Option {
	return func(cfg *kafka.ReaderConfig) {
		cfg.MaxBytes = maxBytes
	}
}

// readMsg 从kafka中读取数据
func (kc *KafkaSyncConsumer) readMsg() {
	// 关闭kafka读取器
	defer kc.reader.Close()
	// channel，1个容量
	sigChan := make(chan os.Signal, 1)
	// syscall.SIGINT 用户在终端按下Ctrl+C，表示中断程序
	// syscall.SIGTERM 系统或管理员发送，表示程序终止
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	// 无限循环
	for {
		select {
		// 监听到终止信号，退出循环
		case <-sigChan:
			zap.L().Info("接收到终止信号，停止消费")
			return
		default:
			// 读取消息
			ctx, cancel := context.WithTimeout(context.Background(), kc.timeout)
			// 从kafka中读取消息
			m, err := kc.reader.ReadMessage(ctx)
			if err != nil {
				if err == context.DeadlineExceeded {
					zap.L().Warn("读取消息超时")
					continue
				}
				zap.L().Error("kafka读取数据失败", zap.Error(err))
				continue
			}
			fmt.Printf("消息在offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
			cancel()
		}
	}
}
