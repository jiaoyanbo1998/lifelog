package ioc

import (
	"lifelog-grpc/event/lifeLogEvent"
	"lifelog-grpc/event/commentEvent"
)

// InitConsumers 初始化单个消费的，消费者组
func InitConsumers(c *lifeLogEvent.ReadEventConsumer) []lifeLogEvent.Consumer {
	// 创建一个单个消费的，消费者组
	consumers := []lifeLogEvent.Consumer{c}
	return consumers
}

// InitBatchConsumers 初始化批量消费的，消费者组
func InitBatchConsumers(c *lifeLogEvent.ReadEventBatchConsumer,
	c1 *commentEvent.CommentEventBatchConsumer) []lifeLogEvent.Consumer {
	// 创建一个批量消费的，消费者组
	consumers := []lifeLogEvent.Consumer{c, c1}
	return consumers
}
