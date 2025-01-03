//go:build wireinject

package main

import (
	"github.com/google/wire"
	"lifelog-grpc/code/grpc"
	"lifelog-grpc/code/ioc"
	"lifelog-grpc/code/repository"
	"lifelog-grpc/code/repository/cache"
	"lifelog-grpc/code/service"
)

// codeSet 注入
var codeSet = wire.NewSet(
	service.NewCodeService,
	repository.NewCodeRepository,
	cache.NewCodeCache,
)

var third = wire.NewSet(
	ioc.InitRedis,
	ioc.InitSms,
	ioc.InitLogger,
	ioc.InitBloomFilter,
)

// InitCodeServiceGRPCService 初始化CodeServiceGRPCService
func InitCodeServiceGRPCService() *grpc.CodeServiceGRPCService {
	wire.Build(
		third,
		grpc.NewCodeServiceGRPCService,
		codeSet,
	)
	return new(grpc.CodeServiceGRPCService)
}
