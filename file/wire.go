//go:build wireinject

package main

import (
	"github.com/google/wire"
	"lifelog-grpc/file/grpc"
	"lifelog-grpc/file/ioc"
	"lifelog-grpc/file/repository"
	"lifelog-grpc/file/repository/dao"
	"lifelog-grpc/file/service"
)

var fileSet = wire.NewSet(
	service.NewFileService,
	repository.NewFileService,
	dao.NewFileService,
)

var thirdSet = wire.NewSet(
	ioc.InitMinio,
	ioc.InitLogger,
	ioc.GetMysql,
)

func InitFileServiceGRPCService() *grpc.FileServiceGRPCService {
	wire.Build(
		fileSet,
		thirdSet,
		grpc.NewFileServiceGRPCService,
	)
	return new(grpc.FileServiceGRPCService)
}
