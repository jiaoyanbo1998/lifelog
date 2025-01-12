package main

import (
	"github.com/spf13/viper"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	feedv1 "lifelog-grpc/api/proto/gen/feed"
	"strconv"
)

func main() {
	// 1.初始化配置文件
	initViperDevelopment()
	// 2.获取配置文件信息
	type feedConfig struct {
		Addr []string `yaml:"addr"`
		Port int      `yaml:"port"`
	}
	c := feedConfig{
		Addr: []string{"localhost:12379"},
		Port: 8077,
	}
	err := viper.UnmarshalKey("feedEtcd", &c)
	if err != nil {
		logx.Error("配置文件读取失败,feed/main")
	}
	// 3.grpc服务端配置文件
	serverConf := zrpc.RpcServerConf{
		ListenOn: ":" + strconv.Itoa(c.Port),
		Etcd: discov.EtcdConf{
			Hosts: c.Addr,
			Key:   "service/feed",
		},
	}
	// 4.初始化feed服务的grpc服务端
	app := InitFeedServiceGRPCService()
	// 5.创建grpc服务端
	server := zrpc.MustNewServer(serverConf, func(server *grpc.Server) {
		// 6.将feed服务注册到grpc服务端
		feedv1.RegisterFeedServiceServer(server, app.feedServiceGRPCService)
	})
	// 7.结束grpc服务端
	defer server.Stop()
	// 8.启动消费者
	err = app.consumer.Start()
	if err != nil {
		logx.Error("消费者启动失败")
	}
	// 9.启动grpc服务端
	logx.Info("正在启动feed服务...，监听的端口为：", strconv.Itoa(c.Port))
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
