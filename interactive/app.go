package main

import (
	"lifelog-grpc/interactive/event/likedEvent"
	"lifelog-grpc/interactive/grpc"
)

type App struct {
	interactiveServiceGRPCService *grpc.InteractiveServiceGRPCService
	asyncLikedEventConsumer       *likedEvent.AsyncLikedEventConsumer
}
