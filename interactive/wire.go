//go:build wireinject

package main

import (
	"github.com/google/wire"
	"lifelog-grpc/interactive/event/feed"
	"lifelog-grpc/interactive/event/likedEvent"
	"lifelog-grpc/interactive/grpc"
	"lifelog-grpc/interactive/ioc"
	"lifelog-grpc/interactive/repository"
	"lifelog-grpc/interactive/repository/cache"
	"lifelog-grpc/interactive/repository/dao"
	"lifelog-grpc/interactive/service"
)

// interactiveSet 注入
var interactiveSet = wire.NewSet(
	service.NewInteractiveService,
	repository.NewInteractiveRepository,
	cache.NewInteractiveCache,
	dao.NewInteractiveDao,
)

var third = wire.NewSet(
	ioc.InitRedis,
	ioc.GetMysql,
	ioc.InitLogger,
	ioc.InitCollectServiceGRPCClient,
	ioc.InitSaramaSyncProducer,
	ioc.InitSaramaKafka,
	ioc.InitFeedServiceGRPCClient,
	ioc.InitLifeLogServiceCRPCClient,
	feed.NewSyncProducer,
	likedEvent.NewAsyncLikedEventConsumer,
)

// InitInteractiveServiceGRPCService 初始化InitInteractiveServiceGRPCService
func InitInteractiveServiceGRPCService() *App {
	wire.Build(
		third,
		grpc.NewCodeServiceGRPCService,
		interactiveSet,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
