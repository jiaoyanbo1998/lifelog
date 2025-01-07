//go:build wireinject

package main

import (
	"github.com/google/wire"
	"lifelog-grpc/comment/grpc"
	"lifelog-grpc/comment/ioc"
	"lifelog-grpc/comment/repository"
	"lifelog-grpc/comment/repository/dao"
	"lifelog-grpc/comment/service"
)

var commentSet = wire.NewSet(
	service.NewCommentService,
	repository.NewCommentRepository,
	dao.NewCommentDaoGorm,
)

var third = wire.NewSet(
	ioc.InitLogger,
	ioc.GetMysql,
	ioc.InitKafkaProducer,
	ioc.NewCommentConsumer,
)

func InitCommentServiceGRPCService() *App {
	wire.Build(
		commentSet,
		third,
		grpc.NewCommentServiceGRPCService,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
