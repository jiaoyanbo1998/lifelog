//go:build wireinject

package main

import (
	"github.com/google/wire"
	"lifelog-grpc/lifeLog/grpc"
	"lifelog-grpc/lifeLog/repository"
	"lifelog-grpc/lifeLog/repository/cache"
	"lifelog-grpc/lifeLog/repository/dao"
	"lifelog-grpc/lifeLog/service"
	"lifelog-grpc/lifelog/ioc"
)

var thirdSet = wire.NewSet(
	ioc.InitRedis,
	ioc.GetMysql,
	ioc.InitLogger,
	ioc.InitGoCache,
)

var lifelogSet = wire.NewSet(
	service.NewLifeLogService,
	repository.NewLifeLogRepository,
	dao.NewLifeLogDao,
	cache.NewLifeLogRedisCache,
	cache.NewLocalCache,
)

func InitLifeLogServiceGRPCService() *grpc.LifeLogServiceGRPCService {
	wire.Build(
		thirdSet,
		lifelogSet,
		grpc.NewLifeLogServiceGRPCService,
	)
	return new(grpc.LifeLogServiceGRPCService)
}
