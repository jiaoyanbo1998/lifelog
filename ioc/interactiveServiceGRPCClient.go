package ioc

import (
	"github.com/spf13/viper"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
	interactivev1 "lifelog-grpc/api/proto/gen/api/proto/interactive/v1"
	"lifelog-grpc/pkg/loggerx"
)

func InitInteractiveServiceGRPCClient(logger loggerx.Logger) interactivev1.InteractiveServiceClient {
	// 1.解析配置文件
	type config struct {
		Addr []string `yaml:"addr"`
	}
	cfg := config{
		Addr: []string{"localhost:12379"},
	}
	err := viper.UnmarshalKey("interactiveEtcd", &cfg)
	if err != nil {
		logger.Error("解析配置文件失败", loggerx.Error(err),
			loggerx.String("method:", "InitInteractiveServiceGRPCClient"))
		panic(err)
	}
	// 2.grpc客户端配置（go-zero框架）
	clientConf := zrpc.RpcClientConf{
		// 服务发现
		Etcd: discov.EtcdConf{
			Hosts: cfg.Addr,              // etcd地址
			Key:   "service/interactive", // 服务的唯一标识
		},
	}
	// 3.创建grpc客户端
	client := zrpc.MustNewClient(clientConf)
	// 4.创建grpc客户端连接
	conn := client.Conn()
	// 5.创建interactiveService的grpc客户端
	return interactivev1.NewInteractiveServiceClient(conn)
}
