//go:build wireinject

package main

import (
	"github.com/google/wire"
	"lifelog-grpc/user/grpc"
	"lifelog-grpc/user/ioc"
	"lifelog-grpc/user/repository"
	"lifelog-grpc/user/repository/cache"
	"lifelog-grpc/user/repository/dao"
	"lifelog-grpc/user/service"
)

// 用户模块
var userSet = wire.NewSet(
	service.NewUserService,
	repository.NewUserRepository,
	dao.NewUserDao,
	cache.NewUserCache,
)

// 第三方模块
var thirdSet = wire.NewSet(
	ioc.InitMysql,
	ioc.InitLogger,
	ioc.InitRedis,
	ioc.InitCodeServiceGRPCClient,
)

func InitUserServiceGRPCServer() *grpc.UserServiceGRPCService {
	wire.Build(
		userSet,
		thirdSet,
		grpc.NewUserServiceGRPCService,
	)
	return new(grpc.UserServiceGRPCService)
}
