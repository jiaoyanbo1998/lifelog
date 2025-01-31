//go:build wireinject

package main

import (
	"github.com/google/wire"
	"lifelog-grpc/event/interactiveEvent"
	"lifelog-grpc/event/lifeLogEvent"
	"lifelog-grpc/ioc"
	"lifelog-grpc/pkg/miniox"
	rankingRepository "lifelog-grpc/ranking/repository"
	rankingCache "lifelog-grpc/ranking/repository/cache"
	rankingService "lifelog-grpc/ranking/service"
	"lifelog-grpc/web"
)

// userSet user模块的依赖注入
var userSet = wire.NewSet(
	web.NewUserHandler,
	ioc.InitUserServiceGRPCClient,
	lifeLogEvent.NewSyncProducer,
)

// codeSet code模块的依赖注入
var codeSet = wire.NewSet(
	web.NewCodeHandler,
	ioc.InitCodeServiceGRPCClient,
)

// JwtSet 初始化jwt模块
var JwtSet = wire.NewSet(
	web.NewJWTHandler,
)

// LifeLog模块
var lifeLogSet = wire.NewSet(
	web.NewLifeLogHandler,
	ioc.InitLifeLogServiceCRPCClient,
	ioc.InitSaramaKafka,              // sarama.Client
	ioc.InitSaramaSyncProducer,       // sarama.SyncProducer
	interactiveEvent.NewSyncProducer, // 参数sarama.AsyncProducer
)

// collectClipSet collectClip模块的依赖注入
var collectClipSet = wire.NewSet(
	web.NewCollectHandler,
	ioc.InitCollectServiceGRPCClient,
)

// interactiveSet interactive模块的依赖注入
var interactiveSet = wire.NewSet(
	web.NewInteractiveHandler,
	ioc.InitInteractiveServiceGRPCClient,
)

// commentSet 评论
var commentSet = wire.NewSet(
	web.NewCommentHandler,
	ioc.InitCommentServiceGRPCClient,
)

// fileSet 文件
var fileSet = wire.NewSet(
	ioc.InitMinio,
	miniox.NewFileHandler,
	web.NewFilesHandler,
	ioc.InitFilesServiceGRPCClient,
)

// feedSet feed流
var feedSet = wire.NewSet(
	web.NewFeedHandler,
	ioc.InitFeedServiceGRPCClient,
)

// rankingSet ranking模块的依赖注入
var rankingSet = wire.NewSet(
	rankingService.NewRankingService,
	rankingRepository.NewRankingRepository,
	rankingCache.NewRankingCacheRedis,
)

// rankingJobCronSet 热榜定时任务的依赖注入
var rankingJobCronSet = wire.NewSet(
	ioc.InitRankingJob,
	ioc.InitCronRankingJob,
)

func InitApp() *App {
	wire.Build(
		// 初始化web服务器
		ioc.InitGin,

		// 初始化日志记录器
		ioc.InitLogger,
		// 初始化web中间件
		ioc.InitMiddlewares,
		// 初始化redis
		ioc.InitRedis,

		// 初始化user模块
		userSet,

		// 初始化lifeLog模块
		lifeLogSet,

		// 初始化jwt模块
		JwtSet,

		// 初始化验证码code模块
		codeSet,

		// 初始化互动信息interactive模块
		interactiveSet,

		// 初始化收藏夹collectClip模块
		collectClipSet,

		// 初始化排行榜ranking模块
		rankingSet,

		// 初始化热榜定时任务
		rankingJobCronSet,

		// 初始化评论comment模块
		commentSet,

		// 初始化文件file模块
		fileSet,

		// 初始化feed流模块
		feedSet,

		// 结构体自动填充
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
