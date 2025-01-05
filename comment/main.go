package main

import (
	"github.com/spf13/viper"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	commentv1 "lifelog-grpc/api/proto/gen/api/proto/comment/v1"
	"strconv"
)

func main() {
	// 1.读取配置文件
	type config struct {
		Port int      `yaml:"port"`
		Addr []string `yaml:"addr"`
	}
	cfg := config{
		Port: 8075,
		Addr: []string{"localhost:12379"},
	}
	err := viper.UnmarshalKey("commentEtcd", &cfg)
	if err != nil {
		panic(err)
	}
	// 2.创建一个grpc服务配置（go-zero框架）
	serverConf := zrpc.RpcServerConf{
		// grpc监听的端口
		ListenOn: ":" + strconv.Itoa(cfg.Port),
		// 服务注册
		Etcd: discov.EtcdConf{
			Hosts: cfg.Addr,          // etcd地址
			Key:   "service/comment", // 服务的唯一标识
		},
	}
	// 3.创建commentService的grpc服务器
	commentServiceGRPCService := InitCommentServiceGRPCService()
	// 4.创建一个grpc服务（go-zero框架）
	server := zrpc.MustNewServer(serverConf, func(server *grpc.Server) {
		// 5.将commentService注册到grpc服务中
		commentv1.RegisterCommentServiceServer(server, commentServiceGRPCService)
	})
	// 6.关闭grpc服务
	defer server.Stop()
	// 7.启动grpc服务
	logx.Info("正在启动commentService的grpc服务...")
	server.Start()
}
