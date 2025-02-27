// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/google/wire"
	"lifelog-grpc/event/interactiveEvent"
	"lifelog-grpc/event/lifeLogEvent"
	"lifelog-grpc/ioc"
	"lifelog-grpc/pkg/miniox"
	"lifelog-grpc/ranking/repository"
	"lifelog-grpc/ranking/repository/cache"
	"lifelog-grpc/ranking/service"
	"lifelog-grpc/web"
)

// Injectors from wire.go:

func InitApp() *App {
	userServiceClient := ioc.InitUserServiceGRPCClient()
	logger := ioc.InitLogger()
	jwtHandler := web.NewJWTHandler(logger)
	minioClient := ioc.InitMinio()
	fileHandler := miniox.NewFileHandler(minioClient)
	client := ioc.InitSaramaKafka(logger)
	syncProducer := ioc.InitSaramaSyncProducer(client)
	lifeLogEventSyncProducer := lifeLogEvent.NewSyncProducer(syncProducer, logger)
	userHandler := web.NewUserHandler(userServiceClient, logger, jwtHandler, fileHandler, lifeLogEventSyncProducer)
	cmdable := ioc.InitRedis()
	v := ioc.InitMiddlewares(logger, cmdable)
	lifeLogServiceClient := ioc.InitLifeLogServiceCRPCClient()
	interactiveServiceClient := ioc.InitInteractiveServiceGRPCClient(logger)
	interactiveEventSyncProducer := interactiveEvent.NewSyncProducer(syncProducer, logger)
	lifeLogHandler := web.NewLifeLogHandler(logger, lifeLogServiceClient, interactiveServiceClient, interactiveEventSyncProducer)
	collectServiceClient := ioc.InitCollectServiceGRPCClient(logger)
	collectHandler := web.NewCollectHandler(collectServiceClient, logger)
	commentServiceClient := ioc.InitCommentServiceGRPCClient(logger)
	commentHandler := web.NewCommentHandler(logger, commentServiceClient)
	codeServiceClient := ioc.InitCodeServiceGRPCClient(logger)
	codeHandler := web.NewCodeHandler(logger, codeServiceClient)
	rankingCache := cache.NewRankingCacheRedis(cmdable)
	rankingRepository := repository.NewRankingRepository(rankingCache)
	rankingService := service.NewRankingService(rankingRepository)
	job := ioc.InitRankingJob(rankingService, logger, cmdable)
	interactiveHandler := web.NewInteractiveHandler(logger, interactiveServiceClient, job)
	filesServiceClient := ioc.InitFilesServiceGRPCClient()
	filesHandler := web.NewFilesHandler(filesServiceClient, fileHandler, logger)
	feedServiceClient := ioc.InitFeedServiceGRPCClient()
	feedHandler := web.NewFeedHandler(feedServiceClient, logger)
	engine := ioc.InitGin(userHandler, v, lifeLogHandler, collectHandler, commentHandler, codeHandler, interactiveHandler, filesHandler, feedHandler)
	cron := ioc.InitCronRankingJob(logger, job)
	app := &App{
		server: engine,
		cron:   cron,
	}
	return app
}

// wire.go:

// userSet user模块的依赖注入
var userSet = wire.NewSet(web.NewUserHandler, ioc.InitUserServiceGRPCClient, lifeLogEvent.NewSyncProducer)

// codeSet code模块的依赖注入
var codeSet = wire.NewSet(web.NewCodeHandler, ioc.InitCodeServiceGRPCClient)

// JwtSet 初始化jwt模块
var JwtSet = wire.NewSet(web.NewJWTHandler)

// LifeLog模块
var lifeLogSet = wire.NewSet(web.NewLifeLogHandler, ioc.InitLifeLogServiceCRPCClient, ioc.InitSaramaKafka, ioc.InitSaramaSyncProducer, interactiveEvent.NewSyncProducer)

// collectClipSet collectClip模块的依赖注入
var collectClipSet = wire.NewSet(web.NewCollectHandler, ioc.InitCollectServiceGRPCClient)

// interactiveSet interactive模块的依赖注入
var interactiveSet = wire.NewSet(web.NewInteractiveHandler, ioc.InitInteractiveServiceGRPCClient)

// commentSet 评论
var commentSet = wire.NewSet(web.NewCommentHandler, ioc.InitCommentServiceGRPCClient)

// fileSet 文件
var fileSet = wire.NewSet(ioc.InitMinio, miniox.NewFileHandler, web.NewFilesHandler, ioc.InitFilesServiceGRPCClient)

// feedSet feed流
var feedSet = wire.NewSet(web.NewFeedHandler, ioc.InitFeedServiceGRPCClient)

// rankingSet ranking模块的依赖注入
var rankingSet = wire.NewSet(service.NewRankingService, repository.NewRankingRepository, cache.NewRankingCacheRedis)

// rankingJobCronSet 热榜定时任务的依赖注入
var rankingJobCronSet = wire.NewSet(ioc.InitRankingJob, ioc.InitCronRankingJob)
