package ioc

import (
	"github.com/spf13/viper"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	feedv1 "lifelog-grpc/api/proto/gen/feed"
)

func InitFeedServiceGRPCClient() feedv1.FeedServiceClient {
	// 1.解析配置文件
	type feedConfig struct {
		Addr []string `yaml:"addr"`
	}
	c := feedConfig{
		Addr: []string{"localhost:12379"},
	}
	err := viper.UnmarshalKey("feedEtcd", &c)
	if err != nil {
		logx.Error("配置文件读取失败")
	}
	// 2.创建grpc客户端配置
	clientConf := zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: c.Addr,
			Key:   "service/feed",
		},
	}
	// 3.创建grpc客户端
	client := zrpc.MustNewClient(clientConf)
	// 4.创建grpc客户端连接
	conn := client.Conn()
	// 5.创建feed服务的grpc客户端
	feedServiceGRPCClient := feedv1.NewFeedServiceClient(conn)
	return feedServiceGRPCClient
}
