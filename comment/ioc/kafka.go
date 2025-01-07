package ioc

import (
	"github.com/spf13/viper"
	"github.com/zeromicro/go-zero/core/logx"
	"lifelog-grpc/pkg/kafkax"
)

func InitKafka() *kafkax.KafkaProducer {
	type kafkaConfig struct {
		Addr  []string `yaml:"addr"`
		Topic string   `yaml:"topic"`
	}
	var kCfg kafkaConfig
	err := viper.UnmarshalKey("commentKafka", &kCfg)
	if err != nil {
		logx.Errorf("加载配置文件失败：%s", err.Error())
		panic(err)
	}
	producer := kafkax.NewKafkaProducer(kCfg.Addr, kCfg.Topic)
	return producer
}
