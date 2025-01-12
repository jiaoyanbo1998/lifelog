package service

import (
	"context"
	"encoding/json"
	"errors"
	"lifelog-grpc/feed/domain"
	"lifelog-grpc/feed/repository"
)

const (
	LikeEventName = "like_event" // 点赞事件的名称
)

// LikeEventHandler 点赞事件处理器，处理点赞的业务逻辑
type LikeEventHandler struct {
	feedRepository repository.FeedRepository
}

func NewLikeEventHandler(feedRepository repository.FeedRepository) Handler {
	return &LikeEventHandler{
		feedRepository: feedRepository,
	}
}

// CreateFeedEvent 创建点赞事件
// 		CreateFeedEvent 中的 ext 里面至少需要三个 id
// 		liked int64:  被点赞的人
// 		liker int64： 点赞的人
// 		bizId int64: 被点赞的东西
// 		biz: string：业务类型
func (l *LikeEventHandler) CreateFeedEvent(ctx context.Context, feedEvent domain.FeedEvent) error {
	// 当用户点赞了，需要通知被点赞的人，谁给你点赞了
	var likeFeedEvent domain.LikeFeedEvent
	err := json.Unmarshal([]byte(feedEvent.Content), &likeFeedEvent)
	if err != nil {
		return errors.New("代码错误，或业务方传递的参数错误")
	}
	return l.feedRepository.CreatePushEvents(ctx, []domain.FeedEvent{
		{
			UserId:     likeFeedEvent.LikedUserId,
			Type:       LikeEventName,
			CreateTime: feedEvent.CreateTime,
			Content:    feedEvent.Content,
		},
	})
}

func (l *LikeEventHandler) FindFeedEvents(ctx context.Context, userId, createTime, limit int64) ([]domain.FeedEvent, error) {
	return l.feedRepository.FindPushEventsWithType(ctx, LikeEventName, userId, createTime, limit)
}
