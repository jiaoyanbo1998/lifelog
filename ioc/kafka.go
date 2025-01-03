package ioc

import (
	"github.com/IBM/sarama"
)

func InitKafka() sarama.Client {
	// 1.kafka地址
	addr := []string{"localhost:9092", "localhost:9094"}

	// 2.创建Kafka配置
	sfg := sarama.NewConfig()

	// 3.返回成功或失败信息
	// 当消息成功发送到broker，返回成功信息
	sfg.Producer.Return.Successes = true
	// 当消息没有成功发送到broker，返回失败信息
	sfg.Producer.Return.Errors = true

	// 4.指定acks（需要等待多少个副本确认）
	// 0：只要消息发送到服务端，就会返回成功信息
	// sfg.Producer.RequiredAcks = sarama.NoResponse // 0
	// 1：需要写入到主分区
	sfg.Producer.RequiredAcks = sarama.WaitForLocal // 1
	// -1：需要写入到所有ISR（ISR == 能够正常同步主分区数据的副本）
	// sfg.Producer.RequiredAcks = sarama.WaitForAll // -1

	// 5.设置分区
	//		sarama默认使用HashPartitioner（还有：轮询，随机策略）
	// 		计算"消息的key"的哈希值，将相同的key存放到一个分区中
	sfg.Producer.Partitioner = sarama.NewHashPartitioner

	// 6.创建Kafka客户端
	client, err := sarama.NewClient(addr, sfg)
	if err != nil {
		panic(err)
	}
	return client
}
