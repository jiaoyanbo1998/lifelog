//go:build wireinject

package main

import (
	"github.com/google/wire"
	"lifelog-grpc/interactive/grpc"
	"lifelog-grpc/interactive/ioc"
	"lifelog-grpc/interactive/repository"
	"lifelog-grpc/interactive/repository/cache"
	"lifelog-grpc/interactive/repository/dao"
	"lifelog-grpc/interactive/service"
)

// codeSet 注入
var codeSet = wire.NewSet(
	service.NewInteractiveService,
	repository.NewInteractiveRepository,
	cache.NewInteractiveCache,
	dao.NewInteractiveDao,
)

var third = wire.NewSet(
	ioc.InitRedis,
	ioc.GetMysql,
	ioc.InitLogger,
)

// InitInteractiveServiceGRPCService 初始化InitInteractiveServiceGRPCService
func InitInteractiveServiceGRPCService() *grpc.InteractiveServiceGRPCService {
	wire.Build(
		third,
		grpc.NewCodeServiceGRPCService,
		codeSet,
	)
	return new(grpc.InteractiveServiceGRPCService)
}
