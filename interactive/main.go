package main

import (
	"github.com/spf13/viper"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"lifelog-grpc/api/proto/gen/interactive/v1"
	"strconv"
)

func main() {
	// 1.读取配置文件
	initViperDevelopment()
	type config struct {
		Port int      `yaml:"port"`
		Addr []string `yaml:"addr"`
	}
	cfg := config{
		Port: 8073,
		Addr: []string{"localhost:12379"},
	}
	err := viper.UnmarshalKey("interactiveEtcd", &cfg)
	if err != nil {
		panic(err)
	}
	// 2.创建一个grpc服务配置（go-zero框架）
	serverConf := zrpc.RpcServerConf{
		// grpc监听的端口
		ListenOn: ":" + strconv.Itoa(cfg.Port),
		// 服务注册
		Etcd: discov.EtcdConf{
			Hosts: cfg.Addr,              // etcd地址
			Key:   "service/interactive", // 服务的唯一标识
		},
	}
	// 3.创建InteractiveService的grpc服务器
	app := InitInteractiveServiceGRPCService()
	// 4.创建一个grpc服务（go-zero框架）
	server := zrpc.MustNewServer(serverConf, func(server *grpc.Server) {
		// 5.将interactiveService注册到grpc服务中
		interactivev1.RegisterInteractiveServiceServer(server, app.interactiveServiceGRPCService)
	})
	// 6.关闭grpc服务
	defer server.Stop()
	// 7.启动消费者
	err = app.asyncLikedEventConsumer.StartConsumer()
	if err != nil {
		panic(err)
	}
	// 8.启动grpc服务
	logx.Info("正在启动interactiveService的grpc服务...")
	server.Start()
}

func initViperDevelopment() {
	// 配置文件的名字
	viper.SetConfigName("dev")
	// 配置文件的类型
	viper.SetConfigType("yaml")
	// 当前目录下的，config目录
	viper.AddConfigPath("config")
	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
