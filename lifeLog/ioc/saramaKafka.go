package ioc

import (
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	"lifelog-grpc/pkg/loggerx"
)

func InitSaramaKafka(logger loggerx.Logger) sarama.Client {
	// 1.kafka地址
	type config struct {
		Addr []string `yaml:"addr"`
	}
	var cfg config
	err := viper.UnmarshalKey("commentKafka", &cfg)
	if err != nil {
		logger.Error("配置文件解析失败", loggerx.Error(err),
			loggerx.String("method:", "comment:ioc:InitSaramaKafka"))
		panic(err)
	}
	// 2.创建Kafka配置
	sfg := sarama.NewConfig()

	// 3.返回成功或失败信息
	// 当消息成功发送到broker，返回成功信息
	sfg.Producer.Return.Successes = true
	// 当消息没有成功发送到broker，返回失败信息
	sfg.Producer.Return.Errors = true

	// 4.指定acks
	//  0：只要消息发送到服务端，就会返回成功信息
	//  sfg.Producer.RequiredAcks = sarama-kafka.NoResponse // 0
	// 1：需要写入到主分区
	sfg.Producer.RequiredAcks = sarama.WaitForLocal // 1
	//  -1：需要写入到所有ISR（ISR == 能够正常同步主分区数据的副本）
	//  sfg.Producer.RequiredAcks = sarama-kafka.WaitForAll // -1

	// 5.设置分区
	//	sarama默认使用HashPartitioner（还有：轮询，随机）
	// 	计算"消息的key"的哈希值，将相同的key存放到一个分区中
	sfg.Producer.Partitioner = sarama.NewHashPartitioner

	// 6.创建Kafka客户端
	client, err := sarama.NewClient(cfg.Addr, sfg)
	if err != nil {
		logger.Error("创建Kafka客户端失败", loggerx.Error(err),
			loggerx.String("method:", "lifeLog:ioc:InitKafka"))
		panic(err)
	}
	return client
}
