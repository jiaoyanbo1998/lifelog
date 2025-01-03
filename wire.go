//go:build wireinject

package main

import (
	"github.com/google/wire"
	codeRepository "lifelog-grpc/code/repository"
	codeCache "lifelog-grpc/code/repository/cache"
	codeService "lifelog-grpc/code/service"
	collectClipRepository "lifelog-grpc/collectClip/repository"
	collectClipDao "lifelog-grpc/collectClip/repository/dao"
	collectClipService "lifelog-grpc/collectClip/service"
	commentRepository "lifelog-grpc/comment/repository"
	commentDao "lifelog-grpc/comment/repository/dao"
	commentService "lifelog-grpc/comment/service"
	"lifelog-grpc/event/commentEvent"
	"lifelog-grpc/event/lifeLogEvent"
	interactiveRepository "lifelog-grpc/interactive/repository"
	interactiveCache "lifelog-grpc/interactive/repository/cache"
	interactiveDao "lifelog-grpc/interactive/repository/dao"
	interactiveService "lifelog-grpc/interactive/service"
	"lifelog-grpc/ioc"
	lifeLogRepository "lifelog-grpc/lifeLog/repository"
	lifeLogCache "lifelog-grpc/lifeLog/repository/cache"
	lifeLogDao "lifelog-grpc/lifeLog/repository/dao"
	lifeLogService "lifelog-grpc/lifeLog/service"
	rankingRepository "lifelog-grpc/ranking/repository"
	rankingCache "lifelog-grpc/ranking/repository/cache"
	rankingService "lifelog-grpc/ranking/service"
	"lifelog-grpc/user/repository"
	"lifelog-grpc/user/repository/cache"
	"lifelog-grpc/user/repository/dao"
	"lifelog-grpc/user/service"
	"lifelog-grpc/web"
)

// userSet user模块的依赖注入
var userSet = wire.NewSet(
	web.NewUserHandler,
	service.NewUserService,
	repository.NewUserRepository,
	dao.NewUserDao,
	cache.NewUserCache,
	ioc.InitUserServiceGRPCClient,
)

// codeSet code模块的依赖注入
var codeSet = wire.NewSet(
	codeService.NewCodeService,
	codeRepository.NewCodeRepository,
	codeCache.NewCodeCache,
	web.NewCodeHandler,
	ioc.InitCodeServiceGRPCClient,
)

// JwtSet 初始化jwt模块
var JwtSet = wire.NewSet(
	web.NewJWTHandler,
)

// 短信模块
var smsSet = wire.NewSet(
	ioc.InitSms,
)

// interactiveSet interactive模块的依赖注入
var interactiveSet = wire.NewSet(
	interactiveService.NewInteractiveService,
	interactiveRepository.NewInteractiveRepository,
	interactiveDao.NewInteractiveDao,
	interactiveCache.NewInteractiveCache,
)

// LifeLog模块
var lifeLogSet = wire.NewSet(
	web.NewLifeLogHandler,
	lifeLogService.NewLifeLogService,
	lifeLogRepository.NewLifeLogRepository,
	lifeLogDao.NewLifeLogDao,
	lifeLogCache.NewLifeLogRedisCache,
)

// collectClipSet collectClip模块的依赖注入
var collectClipSet = wire.NewSet(
	web.NewCollectClipHandler,
	collectClipService.NewCollectClipService,
	collectClipRepository.NewCollectClipRepository,
	collectClipDao.NewCollectClipDao,
)

// kafkaSet kafka模块的依赖注入
var kafkaSet = wire.NewSet(
	ioc.InitKafka,
	lifeLogEvent.NewReadEventBatchConsumer,
	lifeLogEvent.NewReadEventConsumer,
	lifeLogEvent.NewSaramaSyncProducer,
	commentEvent.NewCommentEventBatchConsumer,
	// ioc.InitBatchConsumers,
	ioc.InitConsumers,
	ioc.InitSyncProducer,
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

// commentSet 评论
var commentSet = wire.NewSet(
	web.NewCommentHandler,
	commentService.NewCommentService,
	commentRepository.NewCommentRepository,
	commentDao.NewCommentDaoGorm,
)

func InitApp() *App {
	wire.Build(
		// 初始化web服务器
		ioc.InitGin,

		// 初始化mysql
		// ioc.InitMysql,
		ioc.GetMysql,
		// 初始化日志记录器
		ioc.InitLogger,
		// 初始化web中间件
		ioc.InitMiddlewares,
		// 初始化redis
		ioc.InitRedis,
		// 初始化布隆过滤器
		ioc.InitBloomFilter,

		// 初始化user模块
		userSet,

		// 初始化lifeLog模块
		lifeLogSet,

		// 初始化jwt模块
		JwtSet,

		// 初始化短信sms模块
		smsSet,

		// 初始化验证码code模块
		codeSet,

		// 初始化互动信息interactive模块
		interactiveSet,

		// 初始化收藏夹collectClip模块
		collectClipSet,

		// 初始化kafka模块
		kafkaSet,

		// 初始化排行榜ranking模块
		rankingSet,

		// 初始化热榜定时任务
		rankingJobCronSet,

		// 初始化评论comment模块
		commentSet,

		// 结构体自动填充
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
