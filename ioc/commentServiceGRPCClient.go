package ioc

import (
	"github.com/spf13/viper"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
	commentv1 "lifelog-grpc/api/proto/gen/api/proto/comment/v1"
	"lifelog-grpc/pkg/loggerx"
)

func InitCommentServiceGRPCClient(logger loggerx.Logger) commentv1.CommentServiceClient {
	// 1.解析配置文件
	type config struct {
		Addr []string `yaml:"addr"`
	}
	cfg := config{
		Addr: []string{"localhost:12379"},
	}
	err := viper.UnmarshalKey("commentEtcd", &cfg)
	if err != nil {
		logger.Error("加载配置文件失败", loggerx.Error(err),
			loggerx.String("method:", "comment:InitCommentServiceGRPCClient"))
	}
	// 2.grpc客户端配置（go-zero框架）
	clientConf := zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: cfg.Addr,
			Key:   "service/comment",
		},
	}
	// 3.创建grpc客户端（go-zero框架）
	client := zrpc.MustNewClient(clientConf)
	// 4.创建grpc客户端连接
	conn := client.Conn()
	// 5.创建commentService的grpc客户端
	return commentv1.NewCommentServiceClient(conn)
}
