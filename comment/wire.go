//go:build wireinject

package main

import (
	"github.com/google/wire"
	"lifelog-grpc/comment/event/sarama-kafka"
	"lifelog-grpc/comment/grpc"
	"lifelog-grpc/comment/ioc"
	"lifelog-grpc/comment/repository"
	"lifelog-grpc/comment/repository/dao"
	"lifelog-grpc/comment/service"
)

var commentSet = wire.NewSet(
	service.NewCommentService,
	repository.NewCommentRepository,
	dao.NewCommentDaoGorm,
)

var thirdSet = wire.NewSet(
	ioc.InitLogger,
	ioc.GetMysql,
)

var kafkaSet = wire.NewSet(
	// sarama的配置
	ioc.InitSaramaKafka,
	saramaKafka.NewSyncProducer,
	saramaKafka.NewAsyncCommentEventConsumer,
	saramaKafka.NewAsyncProducer,
	saramaKafka.NewAsyncBatchCommentEventConsumer,
	ioc.InitSaramaSyncProducer,
	// kafka-go的配置
	// kafkago.InitKafkaProducer,
	// kafkago.NewCommentConsumer,
)

func InitCommentServiceGRPCService() *App {
	wire.Build(
		commentSet,
		thirdSet,
		kafkaSet,
		grpc.NewCommentServiceGRPCService,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
