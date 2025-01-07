// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/google/wire"
	"lifelog-grpc/comment/grpc"
	"lifelog-grpc/comment/ioc"
	"lifelog-grpc/comment/repository"
	"lifelog-grpc/comment/repository/dao"
	"lifelog-grpc/comment/service"
)

// Injectors from wire.go:

func InitCommentServiceGRPCService() *App {
	logger := ioc.InitLogger()
	db := ioc.GetMysql(logger)
	commentDao := dao.NewCommentDaoGorm(db, logger)
	commentRepository := repository.NewCommentRepository(commentDao)
	commentService := service.NewCommentService(commentRepository)
	kafkaProducer := ioc.InitKafkaProducer(logger)
	commentServiceGRPCService := grpc.NewCommentServiceGRPCService(commentService, kafkaProducer, logger)
	commentConsumer := ioc.NewCommentConsumer(commentServiceGRPCService, logger)
	app := &App{
		commentServiceGRPCService: commentServiceGRPCService,
		commentConsumer:           commentConsumer,
	}
	return app
}

// wire.go:

var commentSet = wire.NewSet(service.NewCommentService, repository.NewCommentRepository, dao.NewCommentDaoGorm)

var third = wire.NewSet(ioc.InitLogger, ioc.GetMysql, ioc.InitKafkaProducer, ioc.NewCommentConsumer)
