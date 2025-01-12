package main

import (
	"lifelog-grpc/feed/event"
	"lifelog-grpc/feed/grpc"
)

type App struct {
	feedServiceGRPCService *grpc.FeedServiceGRPCService
	consumer               *event.FeedEventAsyncConsumer
}
