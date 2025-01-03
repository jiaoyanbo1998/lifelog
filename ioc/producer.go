package ioc

import "github.com/IBM/sarama"

// InitSyncProducer 初始化同步生产者
func InitSyncProducer(client sarama.Client) sarama.SyncProducer {
	// 使用sarama-kafka客户端，创建同步生产者
	syncProducer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return syncProducer
}

// InitASyncProducer 初始化异步生产者
func InitASyncProducer(client sarama.Client) sarama.AsyncProducer {
	// 使用sarama-kafka客户端，创建异步生产者
	asyncProducer, err := sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return asyncProducer
}
