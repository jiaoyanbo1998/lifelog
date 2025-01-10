package main

import (
	"github.com/spf13/viper"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	filev1 "lifelog-grpc/api/proto/gen/files/v1"
	"strconv"
)

func main() {
	// 1.解析配置文件
	type config struct {
		Addr []string `yaml:"addr"`
		Port int      `yaml:"port"`
	}
	c := config{
		Addr: []string{"localhost:12379"},
		Port: 8076,
	}
	err := viper.UnmarshalKey("fileEtcd", &c)
	if err != nil {
		logx.Error(err)
		logx.Info("解析配置文件失败")
		logx.Info("method:files:main")
	}
	// 2.grpc服务端配置（go-zero框架）
	serverConf := zrpc.RpcServerConf{
		ListenOn: ":" + strconv.Itoa(c.Port),
		Etcd: discov.EtcdConf{
			Hosts: c.Addr,
			Key:   "service/files",
		},
	}
	// 3.创建file服务的grpc服务器
	fileServiceGRPCService := InitFileServiceGRPCService()
	// 4.创建grpc服务端（go-zero框架）
	server := zrpc.MustNewServer(serverConf, func(server *grpc.Server) {
		// 5.将file服务注册到grpc服务端（go-zero框架）
		filev1.RegisterFilesServiceServer(server, fileServiceGRPCService)
	})
	// 6.关闭grpc服务器
	defer server.Stop()
	// 7.启动grpc服务器
	logx.Info("正在启动file服务的grpc服务器，端口:", c.Port)
	server.Start()
}
