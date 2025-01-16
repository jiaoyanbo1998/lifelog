//go:build wireinject

package main

import (
	"github.com/google/wire"
	"lifelog-grpc/collect/grpc"
	"lifelog-grpc/collect/ioc"
	"lifelog-grpc/collect/repository"
	"lifelog-grpc/collect/repository/cache"
	"lifelog-grpc/collect/repository/dao"
	"lifelog-grpc/collect/service"
)

// collectSet 注入
var collectSet = wire.NewSet(
	repository.NewCollectRepository,
	dao.NewCollectDao,
	service.NewCollectService,
	cache.NewCollectRedisCache,
)

var third = wire.NewSet(
	ioc.InitLogger,
	ioc.GetMysql,
	ioc.InitLifeLogServiceCRPCClient,
	ioc.GetRedis,
)

// InitCollectServiceGRPCService 初始化CollectServiceGRPCService
func InitCollectServiceGRPCService() *grpc.CollectServiceGRPCService {
	wire.Build(
		third,
		grpc.NewCollectServiceGRPCService,
		collectSet,
	)
	return new(grpc.CollectServiceGRPCService)
}
