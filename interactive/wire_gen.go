// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/google/wire"
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
	interactiveServiceGRPCService := grpc.NewCodeServiceGRPCService(interactiveService, logger)
	return interactiveServiceGRPCService
}

// wire.go:

// codeSet 注入
var codeSet = wire.NewSet(service.NewInteractiveService, repository.NewInteractiveRepository, cache.NewInteractiveCache, dao.NewInteractiveDao)

var third = wire.NewSet(ioc.InitRedis, ioc.GetMysql, ioc.InitLogger)
