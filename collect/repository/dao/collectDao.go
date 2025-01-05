package dao

import (
	"context"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"lifelog-grpc/collect/domain"
	"lifelog-grpc/pkg/loggerx"
	"time"
)

type CollectDao interface {
	UpdateCollect(ctx context.Context, collect Collect) error
	InsertCollect(ctx context.Context, collect Collect) error
	DeleteCollectByIds(ctx context.Context, ids []int64, authorId int64) error
	PageQuery(ctx context.Context, id int64, limit int, offset int) ([]domain.CollectDomain, error)
	InsertCollectDetail(ctx context.Context, detail CollectDetail) error
	GetCollectDetailById(ctx context.Context, collectId int64, limit int, offset int, authorId int64) ([]domain.CollectDetailDomain, error)
	DeleteCollectDetail(ctx context.Context, collectId, lifeLogId, authorId int64) error
}

type CollectGormDao struct {
	db     *gorm.DB
	logger loggerx.Logger
}

func NewCollectDao(db *gorm.DB, l loggerx.Logger) CollectDao {
	return &CollectGormDao{
		db:     db,
		logger: l,
	}
}

type Collect struct {
	Id         int64  `gorm:"primaryKey;autoIncrement"`
	Name       string `gorm:"uniqueIndex"`
	Status     uint8
	AuthorId   int64
	CreateTime int64
	UpdateTime int64
}

func (Collect) TableName() string {
	return "tb_collect"
}

type CollectDetail struct {
	Id         int64 `gorm:"primaryKey;autoIncrement"` // 主键
	CollectId  int64
	LifeLogId  int64
	CreateTime int64 // 创建时间
	UpdateTime int64 // 更新时间
	Status     uint8
	AuthorId   int64
}

func (CollectDetail) TableName() string {
	return "tb_collect_detail"
}

// UpdateCollect 更新收藏夹
func (c *CollectGormDao) UpdateCollect(ctx context.Context, collect Collect) error {
	err := c.db.WithContext(ctx).Where("id = ? and author_id = ?",
		collect.Id, collect.AuthorId).Model(&Collect{}).
		Updates(map[string]any{
			"name":        collect.Name,
			"update_time": time.Now().UnixMilli(),
		}).Error
	if err != nil {
		c.logger.Error("收藏夹更新失败", loggerx.Error(err),
			loggerx.String("method:", "CollectGormClipDao:UpdateCollect"))
		return err
	}
	return nil
}

// DeleteCollectDetail 删除收藏夹详情
func (c *CollectGormDao) DeleteCollectDetail(ctx context.Context, collectId, lifeLogId, authorId int64) error {
	return c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.WithContext(ctx).Model(&CollectDetail{}).
			Where("collect_id = ? and life_log_id = ? and author_id = ?",
				collectId, lifeLogId, authorId).Updates(map[string]any{
			"status": 2,
		}).Error
		if err != nil {
			c.logger.Error("收藏夹详情删除失败", loggerx.Error(err),
				loggerx.String("method:", "CollectGormClipDao:DeleteCollectDetail"))
			return err
		}
		return nil
	})
}

// InsertCollect 插入收藏夹
func (c *CollectGormDao) InsertCollect(ctx context.Context, collect Collect) error {
	now := time.Now().UnixMilli()
	collect.CreateTime = now
	collect.UpdateTime = now
	collect.Status = 1
	err := c.db.WithContext(ctx).Create(&collect).Error
	if err != nil {
		c.logger.Error("收藏夹插入失败", loggerx.Error(err),
			loggerx.String("method:", "CollectGormClipDao:InsertCollect"))
		return err
	}
	return nil
}

// DeleteCollectByIds 根据id批量删除收藏夹
func (c *CollectGormDao) DeleteCollectByIds(ctx context.Context, ids []int64, authorId int64) error {
	err := c.db.WithContext(ctx).Where("id in ? and author_id = ?",
		ids, authorId).Delete(&Collect{}).Error
	if err != nil {
		c.logger.Error("收藏夹删除失败", loggerx.Error(err),
			loggerx.String("method:", "CollectGormClipDao:DeleteCollectByIds"))
		return err
	}
	return nil
}

// PageQuery 分页查询收藏夹
func (c *CollectGormDao) PageQuery(ctx context.Context, id int64, limit int, offset int) ([]domain.CollectDomain, error) {
	var collects []Collect
	err := c.db.WithContext(ctx).Where("author_id = ?", id).
		Limit(limit).
		Offset(offset).
		Find(&collects).Error
	if err != nil {
		c.logger.Error("收藏夹分页查询失败", loggerx.Error(err),
			loggerx.String("method:", "CollectGormClipDao:PageQuery"))
		return nil, err
	}
	return c.collectsToDomain(collects), nil
}

// collectsToDomain 将收藏夹转换为领域对象
func (c *CollectGormDao) collectsToDomain(clips []Collect) []domain.CollectDomain {
	dcs := make([]domain.CollectDomain, 0, len(clips))
	for _, cl := range clips {
		dcs = append(dcs, domain.CollectDomain{
			Id:         cl.Id,
			Name:       cl.Name,
			AuthorId:   cl.AuthorId,
			Status:     cl.Status,
			CreateTime: cl.CreateTime,
			UpdateTime: cl.UpdateTime,
		})
	}
	return dcs
}

// InsertCollectDetail 将文章插入收藏夹详情
func (c *CollectGormDao) InsertCollectDetail(ctx context.Context, detail CollectDetail) error {
	now := time.Now().UnixMilli()
	detail.CreateTime = now
	detail.UpdateTime = now
	detail.Status = 1
	// 有记录就更新，没有记录才插入
	return c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "collect_id"},  // 冲突判断字段
				{Name: "life_log_id"}, // 冲突判断字段
				{Name: "author_id"},   // 冲突判断字段
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"status", // 更新字段
				"update_time",
			}),
		}).Create(&detail).Error
		if err != nil {
			c.logger.Error("收藏夹详情更新/插入失败", loggerx.Error(err),
				loggerx.String("method:", "CollectGormClipDao:InsertCollectDetail"))
			return err
		}
		return nil
	})
}

// GetCollectDetailById 根据collectId查询收藏夹详情
func (c *CollectGormDao) GetCollectDetailById(ctx context.Context,
	collectId int64,
	limit int, offset int, authorId int64) ([]domain.CollectDetailDomain, error) {
	var collectDetail []CollectDetail
	err := c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.WithContext(ctx).Where("collect_id = ? and author_id = ?",
			collectId, authorId).
			Limit(limit).
			Offset(offset).
			Find(&collectDetail).Error
		return err
	})
	if err != nil {
		c.logger.Error("收藏夹详情查询失败", loggerx.Error(err),
			loggerx.String("method:", "CollectGormClipDao:GetCollectDetailById"))
		return nil, err
	}
	// 将[]collect转为[]domain.CollectDetailDomain
	var collectDetailDomain []domain.CollectDetailDomain
	copier.Copy(&collectDetailDomain, &collectDetail)
	return collectDetailDomain, nil
}
