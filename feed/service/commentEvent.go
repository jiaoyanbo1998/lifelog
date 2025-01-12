package service

import (
	"context"
	"encoding/json"
	"errors"
	"lifelog-grpc/feed/domain"
	"lifelog-grpc/feed/repository"
)

const (
	CommentEventName = "LifeLog_comment_event"
)

type CommentEventHandler struct {
	feedRepository repository.FeedRepository
}

func NewCommentEventHandler(feedRepository repository.FeedRepository) Handler {
	return &CommentEventHandler{
		feedRepository: feedRepository,
	}
}

// CreateFeedEvent 创建feed流事件
// ext扩展结构体，需要，被评论的用户id
func (f *CommentEventHandler) CreateFeedEvent(ctx context.Context, feedEvent domain.FeedEvent) error {
	// 我给A评论了，要通知A谁给你评论了，因此选择Push事件
	var lifeLogCommentEvent domain.LifeLogCommentEvent
	err := json.Unmarshal([]byte(feedEvent.Content), &lifeLogCommentEvent)
	if err != nil {
		return errors.New("代码错误，或业务方传递的参数错误")
	}
	return f.feedRepository.CreatePushEvents(ctx, []domain.FeedEvent{
		{
			UserId:     lifeLogCommentEvent.CommentedUserId,
			Type:       feedEvent.Type,
			CreateTime: feedEvent.CreateTime,
			Content:    feedEvent.Content,
		},
	})
}

func (f *CommentEventHandler) FindFeedEvents(ctx context.Context, userId, createTime, limit int64) ([]domain.FeedEvent, error) {
	return f.feedRepository.FindPushEventsWithType(ctx, CommentEventName, userId, createTime, limit)
}
