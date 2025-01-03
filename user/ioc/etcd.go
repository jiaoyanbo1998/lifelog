package ioc

import (
	"github.com/spf13/viper"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"lifelog-grpc/pkg/loggerx"
)

func InitEtcd(logger loggerx.Logger) *etcdv3.Client {
	type Config struct {
		Addr []string `yaml:"addr"`
	}
	config := Config{
		Addr: []string{"localhost:12379"},
	}
	err := viper.UnmarshalKey("etcd", &config)
	if err != nil {
		panic(err)
	}
	// 创建etcd服务器
	client, err := etcdv3.New(etcdv3.Config{
		Endpoints: config.Addr,
	})
	if err != nil {
		logger.Error("[grpc] 创建etcd客户端失败", loggerx.Error(err))
		panic(err)
	}
	return client
}
