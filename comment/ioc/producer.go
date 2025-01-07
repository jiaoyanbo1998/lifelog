package ioc

import (
	"github.com/spf13/viper"
	"lifelog-grpc/pkg/kafkax"
	"lifelog-grpc/pkg/loggerx"
)

func InitKafkaProducer(logger loggerx.Logger) *kafkax.KafkaProducer {
	type kafkaConfig struct {
		Addr  []string `yaml:"addr"`
		Topic string   `yaml:"topic"`
	}
	var kCfg kafkaConfig
	err := viper.UnmarshalKey("commentKafka", &kCfg)
	if err != nil {
		logger.Error("加载配置文件失败", loggerx.Error(err))
		panic(err)
	}
	producer := kafkax.NewKafkaProducer(kCfg.Addr, kCfg.Topic)
	return producer
}
