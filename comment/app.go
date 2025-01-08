package main

import (
	saramaKafka "lifelog-grpc/comment/event/sarama-kafka"
	"lifelog-grpc/comment/grpc"
)

type App struct {
	commentServiceGRPCService *grpc.CommentServiceGRPCService
	// sarama的消费者
	commentConsumer *saramaKafka.AsyncCommentEventConsumer
	// kafka-go的消费者
	// commentConsumer *kafkago.CommentConsumer
}
