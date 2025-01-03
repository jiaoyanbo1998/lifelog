package main

import (
	"github.com/spf13/viper"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	userv1 "lifelog-grpc/api/proto/gen/user/v1"
	"strconv"
)

func main() {
	// 1.绑定配置文件
	type GRPCConfig struct {
		Port int      `yaml:"port"`
		Addr []string `yaml:"addr"`
	}
	cfg := GRPCConfig{
		Port: 8070,
		Addr: []string{"localhost:12379"},
	}
	err := viper.UnmarshalKey("userEtcd", &cfg)
	if err != nil {
		panic(err)
	}
	// 2.grpc服务端配置
	serverConf := zrpc.RpcServerConf{
		// grpc服务监听的端口
		ListenOn: ":" + strconv.Itoa(cfg.Port),
		// 服务注册
		Etcd: discov.EtcdConf{
			Hosts: cfg.Addr,               // etcd地址
			Key:   "lifeLog/service/user", // 服务的唯一标识
		},
	}
	// 3.创建userService的grpc服务器
	userServiceGRPCServer := InitUserServiceGRPCServer()
	// 4.创建一个grpc服务（go-zero框架）
	grpcServer := zrpc.MustNewServer(serverConf, func(server *grpc.Server) {
		// 5.将userService注册到grpc服务器中
		userv1.RegisterUserServiceServer(server, userServiceGRPCServer)
	})
	// 6.关闭grpc服务器
	defer grpcServer.Stop()
	// 7.启动grpc服务器
	logx.Info("正在启动userService，服务端口：", cfg.Port, "......")
	grpcServer.Start()
}
