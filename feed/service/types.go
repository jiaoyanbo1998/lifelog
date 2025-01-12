package service

import (
	"context"
	"lifelog-grpc/feed/domain"
)

// FeedService 处理业务公共部分
type FeedService interface {
	CreateFeedEvent(ctx context.Context, feedEvent domain.FeedEvent) error
	FindFeedEvents(ctx context.Context, userId, createTime, limit int64) ([]domain.FeedEvent, error)
}

// Handler 处理具体的业务逻辑，不同的业务有不同的handler实现
type Handler interface {
	CreateFeedEvent(ctx context.Context, extendFields domain.FeedEvent) error
	FindFeedEvents(ctx context.Context, userId, createTime, limit int64) ([]domain.FeedEvent, error)
}
