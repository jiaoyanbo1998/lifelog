//go:build wireinject

package main

import (
	"github.com/google/wire"
	"lifelog-grpc/files/grpc"
	"lifelog-grpc/files/ioc"
	"lifelog-grpc/files/repository"
	"lifelog-grpc/files/repository/dao"
	"lifelog-grpc/files/service"
)

var fileSet = wire.NewSet(
	service.NewFileService,
	repository.NewFileService,
	dao.NewFileDao,
)

var thirdSet = wire.NewSet(
	ioc.InitLogger,
	ioc.GetMysql,
)

func InitFileServiceGRPCService() *grpc.FilesServiceGRPCService {
	wire.Build(
		fileSet,
		thirdSet,
		grpc.NewFileServiceGRPCService,
	)
	return new(grpc.FilesServiceGRPCService)
}
