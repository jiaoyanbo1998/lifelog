package repository

import (
	"context"
	"encoding/json"
	"errors"
	"lifelog-grpc/feed/domain"
	"lifelog-grpc/feed/repository/cache"
	"lifelog-grpc/feed/repository/dao"
)

// FeedRepository  Push 推模型 写扩散，将数据写入粉丝的收件箱
//  			   Pull 拉模型 读扩散，将数据写入自己的发件箱
type FeedRepository interface {
	// 创建多个推事件
	CreatePushEvents(ctx context.Context, feedEvents []domain.FeedEvent) error
	// 创建拉事件
	CreatePullEvent(ctx context.Context, feedEvent domain.FeedEvent) error
	// 获取多个拉事件
	FindPullEvents(ctx context.Context, userIds []int64, createTime, limit int64) ([]domain.FeedEvent, error)
	// 获取多个推事件
	FindPushEvents(ctx context.Context, userId, createTime, limit int64) ([]domain.FeedEvent, error)
	// 获取某个类型的，拉事件
	FindPullEventsWithType(ctx context.Context, typ string, userIds []int64, createTime, limit int64) ([]domain.FeedEvent, error)
	// 获取某个类型的，推事件
	FindPushEventsWithType(ctx context.Context, typ string, userId, createTime, limit int64) ([]domain.FeedEvent, error)
}

type FeedRepositoryV1 struct {
	feedPullDAO dao.FeedPullDAO
	feedPushDAO dao.FeedPushDAO
	feedCache   cache.FeedEventCache
}

func NewFeedRepository(feedPullDAO dao.FeedPullDAO, feedPushDAO dao.FeedPushDAO) FeedRepository {
	return &FeedRepositoryV1{
		feedPullDAO: feedPullDAO,
		feedPushDAO: feedPushDAO,
	}
}

func (f FeedRepositoryV1) CreatePushEvents(ctx context.Context, feedEvents []domain.FeedEvent) error {
	// 将domain.FeedEvent转换为dao.FeedPushEvent
	pushEvents := make([]dao.FeedPushEvent, 0, len(feedEvents))
	for _, feedEvent := range feedEvents {
		pushEvents = append(pushEvents, convertToPushEventDao(feedEvent))
	}
	// 创建Push事件
	return f.feedPushDAO.CreatePushEvents(ctx, pushEvents)
}

// CreatePullEvent 创建拉事件
func (f FeedRepositoryV1) CreatePullEvent(ctx context.Context, feedEvent domain.FeedEvent) error {
	return f.feedPullDAO.CreatePullEvent(ctx, convertToPullEventDao(feedEvent))
}

// FindFeedEvents 查找多个拉事件
func (f FeedRepositoryV1) FindPullEvents(ctx context.Context, userIds []int64, createTime, limit int64) ([]domain.FeedEvent, error) {
	pullEvents, err := f.feedPullDAO.FindPullEventList(ctx, userIds, createTime, limit)
	if err != nil {
		return nil, err
	}
	// 将dao.FeedPullEvent转换为domain.FeedEvent
	feedEvents := make([]domain.FeedEvent, 0, len(pullEvents))
	for _, pullEvent := range pullEvents {
		feedEvents = append(feedEvents, convertToPullEventDomain(pullEvent))
	}
	return feedEvents, nil
}

// FindFeedEvents 查找多个推事件
func (f FeedRepositoryV1) FindPushEvents(ctx context.Context, userId, createTime, limit int64) ([]domain.FeedEvent, error) {
	pushEvents, err := f.feedPushDAO.GetPushEvents(ctx, userId, createTime, limit)
	if err != nil {
		return nil, err
	}
	// 将dao.FeedPushEvent转换为domain.FeedEvent
	feedEvents := make([]domain.FeedEvent, 0, len(pushEvents))
	for _, pushEvent := range pushEvents {
		feedEvents = append(feedEvents, convertToPushEventDomain(pushEvent))
	}
	return feedEvents, nil
}

// FindFeedEventsWithType 查找某个类型的拉事件
func (f FeedRepositoryV1) FindPullEventsWithType(ctx context.Context, typ string, userIds []int64, createTime, limit int64) ([]domain.FeedEvent, error) {
	pullEvents, err := f.feedPullDAO.FindPullEventListWithType(ctx, typ, userIds, createTime, limit)
	if err != nil {
		return nil, err
	}
	// 将dao.FeedPullEvent转换为domain.FeedEvent
	feedEvents := make([]domain.FeedEvent, 0, len(pullEvents))
	for _, pullEvent := range pullEvents {
		feedEvents = append(feedEvents, convertToPullEventDomain(pullEvent))
	}
	return feedEvents, nil
}

// FindFeedEventsWithType 查找某个类型的推事件
func (f FeedRepositoryV1) FindPushEventsWithType(ctx context.Context, typ string, userId, createTime, limit int64) ([]domain.FeedEvent, error) {
	feedPushs, err := f.feedPushDAO.GetPushEventsWithType(ctx, typ, userId, createTime, limit)
	if err != nil {
		return nil, err
	}
	// 将dao.FeedPushEvent转换为domain.FeedEvent
	feedEvents := make([]domain.FeedEvent, 0, len(feedPushs))
	for _, feedPush := range feedPushs {
		feedEvents = append(feedEvents, convertToPushEventDomain(feedPush))
	}
	return feedEvents, nil
}

func (f *FeedRepositoryV1) SetFollowees(ctx context.Context, follower int64, followees []int64) error {
	return f.feedCache.SetFollowees(ctx, follower, followees)
}

func (f *FeedRepositoryV1) GetFollowees(ctx context.Context, follower int64) ([]int64, error) {
	followees, err := f.feedCache.GetFollowees(ctx, follower)
	if errors.Is(err, cache.FolloweesNotFound) {
		return nil, errors.New("没有找到关注列表")
	}
	return followees, err
}

// convertToPushEventDao 将domain.FeedEvent转换为dao.FeedPushEvent
func convertToPushEventDao(feedEvent domain.FeedEvent) dao.FeedPushEvent {
	return dao.FeedPushEvent{
		Id:         feedEvent.ID,
		UserId:     feedEvent.UserId,
		Type:       feedEvent.Type,
		Content:    feedEvent.Content,
		CreateTime: feedEvent.CreateTime,
	}
}

// convertToPullEventDao 将domain.FeedEvent转换为dao.FeedPullEvent
func convertToPullEventDao(feedEvent domain.FeedEvent) dao.FeedPullEvent {
	return dao.FeedPullEvent{
		Id:         feedEvent.ID,
		UserId:     feedEvent.UserId,
		Type:       feedEvent.Type,
		Content:    feedEvent.Content,
		CreateTime: feedEvent.CreateTime,
	}

}

// convertToPushEventDomain 将dao.FeedPushEvent转换为domain.FeedEvent
func convertToPushEventDomain(feedEvent dao.FeedPushEvent) domain.FeedEvent {
	var ext map[string]any
	_ = json.Unmarshal([]byte(feedEvent.Content), &ext)
	return domain.FeedEvent{
		ID:         feedEvent.Id,
		UserId:     feedEvent.UserId,
		Type:       feedEvent.Type,
		CreateTime: feedEvent.CreateTime,
		Content:    feedEvent.Content,
	}
}

// convertToPullEventDomain 将dao.FeedPullEvent转换为domain.FeedEvent
func convertToPullEventDomain(feedEvent dao.FeedPullEvent) domain.FeedEvent {
	var ext map[string]any
	_ = json.Unmarshal([]byte(feedEvent.Content), &ext)
	return domain.FeedEvent{
		ID:         feedEvent.Id,
		UserId:     feedEvent.UserId,
		Type:       feedEvent.Type,
		CreateTime: feedEvent.CreateTime,
		Content:    feedEvent.Content,
	}
}
