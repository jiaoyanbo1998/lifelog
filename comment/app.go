package main

import (
	"lifelog-grpc/comment/grpc"
	"lifelog-grpc/comment/ioc"
)

type App struct {
	commentServiceGRPCService *grpc.CommentServiceGRPCService
	commentConsumer           *ioc.CommentConsumer
}
