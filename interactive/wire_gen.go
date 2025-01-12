// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/google/wire"
	"lifelog-grpc/interactive/event/feed"
	"lifelog-grpc/interactive/grpc"
	"lifelog-grpc/interactive/ioc"
	"lifelog-grpc/interactive/repository"
	"lifelog-grpc/interactive/repository/cache"
	"lifelog-grpc/interactive/repository/dao"
	"lifelog-grpc/interactive/service"
)

// Injectors from wire.go:

// InitInteractiveServiceGRPCService 初始化InitInteractiveServiceGRPCService
func InitInteractiveServiceGRPCService() *grpc.InteractiveServiceGRPCService {
	logger := ioc.InitLogger()
	db := ioc.GetMysql(logger)
	interactiveDao := dao.NewInteractiveDao(db, logger)
	cmdable := ioc.InitRedis()
	interactiveCache := cache.NewInteractiveCache(cmdable, logger)
	interactiveRepository := repository.NewInteractiveRepository(interactiveDao, interactiveCache)
	interactiveService := service.NewInteractiveService(interactiveRepository)
	collectServiceClient := ioc.InitCollectServiceGRPCClient(logger)
	feedServiceClient := ioc.InitFeedServiceGRPCClient()
	client := ioc.InitSaramaKafka(logger)
	syncProducer := ioc.InitSaramaSyncProducer(client)
	feedSyncProducer := feed.NewSyncProducer(syncProducer, logger)
	lifeLogServiceClient := ioc.InitLifeLogServiceCRPCClient()
	interactiveServiceGRPCService := grpc.NewCodeServiceGRPCService(interactiveService, collectServiceClient, feedServiceClient, feedSyncProducer, lifeLogServiceClient)
	return interactiveServiceGRPCService
}

// wire.go:

// interactiveSet 注入
var interactiveSet = wire.NewSet(service.NewInteractiveService, repository.NewInteractiveRepository, cache.NewInteractiveCache, dao.NewInteractiveDao)

var third = wire.NewSet(ioc.InitRedis, ioc.GetMysql, ioc.InitLogger, ioc.InitCollectServiceGRPCClient, ioc.InitSaramaSyncProducer, ioc.InitSaramaKafka, ioc.InitFeedServiceGRPCClient, ioc.InitLifeLogServiceCRPCClient, feed.NewSyncProducer)
