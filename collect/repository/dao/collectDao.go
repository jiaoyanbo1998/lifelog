package dao

import (
	"context"
	"gorm.io/gorm"
	"lifelog-grpc/collect/domain"
	"lifelog-grpc/pkg/loggerx"
	"time"
)

type CollectDao interface {
	UpdateCollect(ctx context.Context, collect Collect) error
	InsertCollect(ctx context.Context, collect Collect) error
	DeleteCollectByIds(ctx context.Context, ids []int64) error
	PageQuery(ctx context.Context, id int64, limit int, offset int) ([]domain.CollectDomain, error)
	InsertCollectDetail(ctx context.Context, detail CollectDetail) error
	GetCollectDetailById(ctx context.Context, id int64, limit int, offset int, authorId int64) ([]domain.CollectDetailDomain, error)
}

type CollectGormClipDao struct {
	db     *gorm.DB
	logger loggerx.Logger
}

func NewCollectDao(db *gorm.DB, l loggerx.Logger) CollectDao {
	return &CollectGormClipDao{
		db:     db,
		logger: l,
	}
}

type Collect struct {
	Id         int64  `gorm:"primaryKey;autoIncrement"`
	Name       string `gorm:"uniqueIndex"`
	Status     uint8
	UserId     int64
	CreateTime int64
	UpdateTime int64
}

func (Collect) TableName() string {
	return "tb_collect_clip"
}

type Interactive struct {
	Id           int64  `gorm:"primaryKey;autoIncrement"`   // 主键
	Biz          string `gorm:"uniqueIndex:idx_biz_id_biz"` // 业务类型
	BizId        int64  `gorm:"uniqueIndex:idx_biz_id_biz"` // 业务id（文章id）
	ReadCount    int64  // 阅读数
	CollectCount int64  // 收藏数
	LikeCount    int64  // 点赞数
	CreateTime   int64  // 创建时间
	UpdateTime   int64  // 更新时间
}

func (Interactive) TableName() string {
	return "tb_interactive"
}

type CollectDetail struct {
	Id         int64 `gorm:"primaryKey;autoIncrement"` // 主键
	CollectId  int64
	LifeLogId  int64
	CreateTime int64 // 创建时间
	UpdateTime int64 // 更新时间
	Status     uint8
}

func (CollectDetail) TableName() string {
	return "tb_collect_clip_detail"
}

type PublicLifeLog struct {
	Id         int64
	Title      string
	Content    string
	AuthorId   int64
	CreateTime int64
	UpdateTime int64
	Status     uint8
}

func (PublicLifeLog) TableName() string {
	return "tb_publish_lifeLog"
}

type CollectDetailWithLifeLog struct {
	CollectDetail `gorm:"embedded"` // 嵌入CollectDetail结构体
	PublicLifeLog `gorm:"embedded"` // 嵌入PublicLifeLog结构体
}

// UpdateCollect 更新收藏夹
func (c *CollectGormClipDao) UpdateCollect(ctx context.Context, collect Collect) error {
	err := c.db.WithContext(ctx).Where("id = ?", collect.Id).Model(&Collect{}).
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

// InsertCollect 插入收藏夹
func (c *CollectGormClipDao) InsertCollect(ctx context.Context, collect Collect) error {
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
func (c *CollectGormClipDao) DeleteCollectByIds(ctx context.Context, ids []int64) error {
	err := c.db.WithContext(ctx).Where("id in ?", ids).Delete(&Collect{}).Error
	if err != nil {
		c.logger.Error("收藏夹删除失败", loggerx.Error(err),
			loggerx.String("method:", "CollectGormClipDao:DeleteCollectByIds"))
		return err
	}
	return nil
}

// PageQuery 分页查询收藏夹
func (c *CollectGormClipDao) PageQuery(ctx context.Context, id int64, limit int, offset int) ([]domain.CollectDomain, error) {
	var collects []Collect
	err := c.db.WithContext(ctx).Where("user_id = ?", id).
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
func (c *CollectGormClipDao) collectsToDomain(clips []Collect) []domain.CollectDomain {
	dcs := make([]domain.CollectDomain, 0, len(clips))
	for _, cl := range clips {
		dcs = append(dcs, domain.CollectDomain{
			Id:         cl.Id,
			Name:       cl.Name,
			UserId:     cl.UserId,
			Status:     cl.Status,
			CreateTime: cl.CreateTime,
			UpdateTime: cl.UpdateTime,
		})
	}
	return dcs
}

// InsertCollectDetail 将文章插入收藏夹详情
func (c *CollectGormClipDao) InsertCollectDetail(ctx context.Context, detail CollectDetail) error {
	now := time.Now().UnixMilli()
	detail.CreateTime = now
	detail.UpdateTime = now
	detail.Status = 1
	err := c.db.WithContext(ctx).Create(&detail).Error
	if err != nil {
		c.logger.Error("收藏夹详情插入失败", loggerx.Error(err),
			loggerx.String("method:", "CollectGormClipDao:InsertCollectDetail"))
		return err
	}
	return nil
}

// GetCollectDetailById 根据id查询收藏夹详情
func (c *CollectGormClipDao) GetCollectDetailById(ctx context.Context, id int64,
	limit int, offset int, authorId int64) ([]domain.CollectDetailDomain, error) {
	var cwas []CollectDetailWithLifeLog
	// 执行JOIN 查询
	err := c.db.WithContext(ctx).
		Table("tb_collect_clip_detail as clip_detail").
		Select(`
            clip_detail.id as id,
            clip_detail.collect_id as collect_id,
            clip_detail.lifeLog_id as lifeLog_id,
            clip_detail.create_time as create_time,
            clip_detail.update_time as update_time,
            clip_detail.status as status,
            lifeLog.id as lifeLog_id,
            lifeLog.title as title,
            lifeLog.content as content,
            lifeLog.author_id as author_id,
            lifeLog.create_time as lifeLog_create_time,
            lifeLog.update_time as lifeLog_update_time,
            lifeLog.status as lifeLog_status
        `).
		Joins("inner join tb_publish_lifeLog as lifeLog on clip_detail.lifeLog_id = lifeLog.id").
		Where("clip_detail.collect_id = ? and clip_detail.status != ? and "+
			"clip_detail.author_id = ?", id, 2, authorId).
		Limit(limit).
		Offset(offset).
		Scan(&cwas).Error
	if err != nil {
		return nil, err
	}
	// 映射查询结果到 CollectDetailDomain
	var collectDetails []domain.CollectDetailDomain
	for _, cwa := range cwas {
		collectDetail := domain.CollectDetailDomain{
			Id:         cwa.CollectDetail.Id,
			CollectId:  cwa.CollectDetail.CollectId,
			LifeLogId:  cwa.CollectDetail.LifeLogId,
			CreateTime: cwa.CollectDetail.CreateTime,
			UpdateTime: cwa.CollectDetail.UpdateTime,
			Status:     cwa.CollectDetail.Status,
			PublicLifeLogDomain: domain.PublicLifeLogDomain{
				Id:         cwa.PublicLifeLog.Id,
				Title:      cwa.PublicLifeLog.Title,
				Content:    cwa.PublicLifeLog.Content,
				AuthorId:   cwa.PublicLifeLog.AuthorId,
				CreateTime: cwa.PublicLifeLog.CreateTime,
				UpdateTime: cwa.PublicLifeLog.UpdateTime,
				Status:     cwa.PublicLifeLog.Status,
			},
		}
		collectDetails = append(collectDetails, collectDetail)
	}
	return collectDetails, nil
}
