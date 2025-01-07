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
	ioc.InitKafka,
)

func InitCommentServiceGRPCService() *grpc.CommentServiceGRPCService {
	wire.Build(
		commentSet,
		third,
		grpc.NewCommentServiceGRPCService,
	)
	return new(grpc.CommentServiceGRPCService)
}
