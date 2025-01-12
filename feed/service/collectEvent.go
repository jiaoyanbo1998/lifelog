package service

import (
	"context"
	"encoding/json"
	"errors"
	"lifelog-grpc/feed/domain"
	"lifelog-grpc/feed/repository"
)

const (
	CollectEventName = "collect_event" // 点赞事件的名称
)

// CollectEventHandler 收藏事件处理器，处理收藏的业务逻辑
type CollectEventHandler struct {
	feedRepository repository.FeedRepository
}

func NewCollectEventHandler(feedRepository repository.FeedRepository) Handler {
	return &CollectEventHandler{
		feedRepository: feedRepository,
	}
}

// CreateFeedEvent 创建收藏事件
func (l *CollectEventHandler) CreateFeedEvent(ctx context.Context, feedEvent domain.FeedEvent) error {
	// 当用户点赞了，需要通知被点赞的人，谁给你点赞了
	var collectFeedEvent domain.CollectFeedEvent
	err := json.Unmarshal([]byte(feedEvent.Content), &collectFeedEvent)
	if err != nil {
		return errors.New("代码错误，或业务方传递的参数错误")
	}
	return l.feedRepository.CreatePushEvents(ctx, []domain.FeedEvent{
		{
			UserId:     collectFeedEvent.CollectedUserId,
			Type:       CollectEventName,
			CreateTime: feedEvent.CreateTime,
			Content:    feedEvent.Content,
		},
	})
}

func (l *CollectEventHandler) FindFeedEvents(ctx context.Context, userId, createTime, limit int64) ([]domain.FeedEvent, error) {
	return l.feedRepository.FindPushEventsWithType(ctx, CollectEventName, userId, createTime, limit)
}
