package dao

import (
	"context"
	"gorm.io/gorm"
)

type FeedPullDAO interface {
	// CreatePullEvent 创建拉事件
	CreatePullEvent(ctx context.Context, event FeedPullEvent) error
	// FindPullEventList 查找多个拉事件
	FindPullEventList(ctx context.Context, uids []int64, timestamp, limit int64) ([]FeedPullEvent, error)
	// FindPullEventListWithType 根据Type查找多个拉事件
	FindPullEventListWithType(ctx context.Context, typ string, uids []int64, timestamp, limit int64) ([]FeedPullEvent, error)
}

type FeedPullGormDAO struct {
	db *gorm.DB
}

func NewFeedPullGormDAO(db *gorm.DB) FeedPullDAO {
	return &FeedPullGormDAO{
		db: db,
	}
}

type FeedPullEvent struct {
	Id      int64
	UserId  int64
	Type    string
	Content string
	// 发生时间
	CreateTime int64
}

func (FeedPullEvent) TableName() string {
	return "tb_feed"
}

func (f *FeedPullGormDAO) FindPullEventListWithType(ctx context.Context,
	typ string, userIds []int64, createTime, limit int64) ([]FeedPullEvent, error) {
	var feedPullEvents []FeedPullEvent
	err := f.db.WithContext(ctx).
		Where("user_id in ?", userIds). // 多个用户id
		Where("create_time < ?", createTime). // 获取小于createTime的事件
		Where("type = ?", typ). // 获取某个类型的事件
		Order("create_time, id  desc"). // 倒序，从大到小排序，最新的数据在前面
		Limit(int(limit)).
		Find(&feedPullEvents).Error
	return feedPullEvents, err
}

func (f *FeedPullGormDAO) CreatePullEvent(ctx context.Context, event FeedPullEvent) error {
	return f.db.WithContext(ctx).Create(&event).Error
}

func (f *FeedPullGormDAO) FindPullEventList(ctx context.Context, userIds []int64,
	createTime, limit int64) ([]FeedPullEvent, error) {
	var feedPullEvents []FeedPullEvent
	err := f.db.WithContext(ctx).
		Where("user_id in ?", userIds).
		Where("create_time < ?", createTime).
		Order("create_time, id  desc").
		Limit(int(limit)).
		Find(&feedPullEvents).Error
	return feedPullEvents, err
}
