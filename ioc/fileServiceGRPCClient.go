package ioc

import (
	"github.com/spf13/viper"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	filev1 "lifelog-grpc/api/proto/gen/file/v1"
)

func InitFileServiceGRPCClient() filev1.FileServiceClient {
	// 1.解析配置文件
	type config struct {
		Addr []string `yaml:"addr"`
	}
	c := config{
		Addr: []string{"localhost:12379"},
	}
	err := viper.UnmarshalKey("fileEtcd", &c)
	if err != nil {
		logx.Error(err)
		logx.Info("读取配置文件失败", "method:ioc:InitFileServiceGRPCClient")
	}
	// grpc客户端配置
	clientConf := zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: c.Addr,
			Key:   "service/file",
		},
	}
	// 创建grpc客户端
	client := zrpc.MustNewClient(clientConf)
	// 创建grpc连接
	conn := client.Conn()
	// 创建file服务的grpc服务客户端
	return filev1.NewFileServiceClient(conn)
}
