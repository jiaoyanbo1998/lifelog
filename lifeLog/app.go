package main

import (
	"lifelog-grpc/lifeLog/event"
	"lifelog-grpc/lifeLog/grpc"
)

type App struct {
	lifeLogServiceGRPCService *grpc.LifeLogServiceGRPCService
	asyncLifeLogEventConsumer *event.AsyncLifeLogEventConsumer
}
