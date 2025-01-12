package service

import (
	"context"
	"encoding/json"
	"errors"
	"lifelog-grpc/feed/domain"
	"lifelog-grpc/feed/repository"
)

const (
	ReadEventName = "read_event" // 点赞事件的名称
)

// ReadEventHandler 阅读事件处理器，处理阅读的业务逻辑
type ReadEventHandler struct {
	feedRepository repository.FeedRepository
}

func NewReadEventHandler(feedRepository repository.FeedRepository) Handler {
	return &ReadEventHandler{
		feedRepository: feedRepository,
	}
}

// CreateFeedEvent 创建收藏事件
func (l *ReadEventHandler) CreateFeedEvent(ctx context.Context, feedEvent domain.FeedEvent) error {
	// 当用户点赞了，需要通知被点赞的人，谁给你点赞了
	var readFeedEvent domain.ReadFeedEvent
	err := json.Unmarshal([]byte(feedEvent.Content), &readFeedEvent)
	if err != nil {
		return errors.New("代码错误，或业务方传递的参数错误")
	}
	return l.feedRepository.CreatePushEvents(ctx, []domain.FeedEvent{
		{
			UserId:     readFeedEvent.ReadedUserId,
			Type:       ReadEventName,
			CreateTime: feedEvent.CreateTime,
			Content:    feedEvent.Content,
		},
	})
}

func (l *ReadEventHandler) FindFeedEvents(ctx context.Context, userId, createTime, limit int64) ([]domain.FeedEvent, error) {
	return l.feedRepository.FindPushEventsWithType(ctx, ReadEventName, userId, createTime, limit)
}
