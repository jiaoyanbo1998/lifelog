package ioc

import (
	"github.com/spf13/viper"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
	"lifelog-grpc/api/proto/gen/lifelog/v1"
)

func InitLifeLogServiceCRPCClient() lifelogv1.LifeLogServiceClient {
	// 1.解析配置文件
	type config struct {
		Addr []string `yaml:"addr"`
	}
	cfg := config{
		Addr: []string{"localhost:12379"},
	}
	err := viper.UnmarshalKey("lifeLogEtcd", &cfg)
	if err != nil {
		panic(err)
	}
	// 2.grpc客户端配置（go-zero）
	clientConf := zrpc.RpcClientConf{
		// 服务发现
		Etcd: discov.EtcdConf{
			Hosts: cfg.Addr,
			Key:   "service/lifelog",
		},
	}
	// 3.创建grpc客户端（go-zero）
	client := zrpc.MustNewClient(clientConf)
	// 4.创建grpc连接
	conn := client.Conn()
	// 5.创建lifelog服务的grpc客户端
	return lifelogv1.NewLifeLogServiceClient(conn)
}
