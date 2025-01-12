package service

import (
	"context"
	"encoding/json"
	"errors"
	"lifelog-grpc/feed/domain"
	"lifelog-grpc/feed/repository"
)

const (
	FollowEventName = "follow_event"
)

type FollowEventHandler struct {
	feedRepository repository.FeedRepository
}

func NewFollowEventHandler(feedRepository repository.FeedRepository) Handler {
	return &FollowEventHandler{
		feedRepository: feedRepository,
	}
}

// CreateFeedEvent 创建跟随方式
// 		如果 A 关注了 B，那么
// 		follower 就是 A
// 		followee 就是 B
func (f *FollowEventHandler) CreateFeedEvent(ctx context.Context, feedEvent domain.FeedEvent) error {
	// 当用户关注了，需要通知被关注的人，谁给你关注了
	var followFeedEvent domain.FollowFeedEvent
	err := json.Unmarshal([]byte(feedEvent.Content), &followFeedEvent)
	if err != nil {
		return errors.New("代码错误，或业务方传递的参数错误")
	}
	return f.feedRepository.CreatePushEvents(ctx, []domain.FeedEvent{
		{
			UserId:     followFeedEvent.FollowedUserId,
			Type:       feedEvent.Type,
			CreateTime: feedEvent.CreateTime,
			Content:    feedEvent.Content,
		},
	})
}

func (f *FollowEventHandler) FindFeedEvents(ctx context.Context, userId, createTime, limit int64) ([]domain.FeedEvent, error) {
	return f.feedRepository.FindPushEventsWithType(ctx, FollowEventName, userId, createTime, limit)
}
