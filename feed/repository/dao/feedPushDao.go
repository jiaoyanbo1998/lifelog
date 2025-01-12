package dao

import (
	"context"
	"gorm.io/gorm"
)

type FeedPushDAO interface {
	// CreatePushEvents 创建多个推事件
	CreatePushEvents(ctx context.Context, events []FeedPushEvent) error
	// GetPushEvents 获取多个推事件
	GetPushEvents(ctx context.Context, uid int64, timestamp, limit int64) ([]FeedPushEvent, error)
	// GetPushEventsWithType 根据Type获取多个推事件
	GetPushEventsWithType(ctx context.Context, typ string, uid int64, timestamp, limit int64) ([]FeedPushEvent, error)
}

type FeedPushGormDAO struct {
	db *gorm.DB
}

func NewFeedPushGormDAO(db *gorm.DB) FeedPushDAO {
	return &FeedPushGormDAO{
		db: db,
	}
}

type FeedPushEvent struct {
	Id      int64
	UserId  int64
	Type    string
	Content string
	// 发生时间
	CreateTime int64
}

func (FeedPushEvent) TableName() string {
	return "tb_feed"
}

func (f *FeedPushGormDAO) GetPushEventsWithType(ctx context.Context, typ string,
	userId int64, createTime, limit int64) ([]FeedPushEvent, error) {
	var events []FeedPushEvent
	err := f.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Where("create_time < ?", createTime).
		Where("type = ?", typ).
		Order("create_time, id  desc").
		Limit(int(limit)).
		Find(&events).Error
	return events, err
}

func (f *FeedPushGormDAO) CreatePushEvents(ctx context.Context, feedPushEvents []FeedPushEvent) error {
	return f.db.WithContext(ctx).Create(&feedPushEvents).Error
}

func (f *FeedPushGormDAO) GetPushEvents(ctx context.Context, userId int64, createTime, limit int64) ([]FeedPushEvent, error) {
	var events []FeedPushEvent
	err := f.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Where("create_time < ?", createTime).
		Order("create_time, id  desc").
		Limit(int(limit)).
		Find(&events).Error
	return events, err
}
