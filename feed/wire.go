//go:build wireinject

package main

import (
	"github.com/google/wire"
	"lifelog-grpc/feed/event"
	"lifelog-grpc/feed/grpc"
	"lifelog-grpc/feed/ioc"
	"lifelog-grpc/feed/repository"
	"lifelog-grpc/feed/repository/dao"
	"lifelog-grpc/feed/service"
)

var feedSet = wire.NewSet(
	service.NewFeedService,
	repository.NewFeedRepository,
	dao.NewFeedPullGormDAO,
	dao.NewFeedPushGormDAO,
)

var thirdSet = wire.NewSet(
	ioc.InitRedis,
	ioc.GetMysql,
	ioc.InitLogger,
	ioc.RegisterHandler,
	ioc.InitSaramaKafka,
)

func InitFeedServiceGRPCService() *App {
	wire.Build(
		thirdSet,
		feedSet,
		event.NewFeedEventAsyncConsumer,
		grpc.NewFeedServiceGRPCService,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
