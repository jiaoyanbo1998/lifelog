package ioc

import (
	"github.com/spf13/viper"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
	userv1 "lifelog-grpc/api/proto/gen/user/v1"
)

// InitUserServiceGRPCClient 初始化grpc用户服务客户端，引入etcd服务发现
func InitUserServiceGRPCClient() userv1.UserServiceClient {
	// 1.解析配置文件
	type userConfig struct {
		Addr []string `yaml:"addr"`
	}
	cfg := userConfig{
		Addr: []string{"localhost:12379"},
	}
	err := viper.UnmarshalKey("userEtcd", &cfg)
	if err != nil {
		panic(err)
	}
	// 2.grpc客户端配置（go-zero框架）
	clientConf := zrpc.RpcClientConf{
		// 服务发现
		Etcd: discov.EtcdConf{
			// etcd地址
			Hosts: cfg.Addr,
			// 服务的唯一标识
			Key: "lifeLog/service/user",
		},
	}
	// 3.创建grpc客户端（go-zero框架）
	client := zrpc.MustNewClient(clientConf)
	// 4.创建grpc客户端连接
	conn := client.Conn()
	// 5.创建userService服务的grpc客户端
	return userv1.NewUserServiceClient(conn)
}
