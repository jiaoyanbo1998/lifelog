package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
	"lifelog-grpc/interactive/domain"
	"lifelog-grpc/pkg/loggerx"
)

type InteractiveDao interface {
	InsertReadCount(ctx context.Context, biz string, bizId int64, userId int64) error
	InsertReadInfo(ctx context.Context, biz string, bizId int64, userId int64) error
	InsertLikeCount(ctx context.Context, biz string, bizId int64, userId int64) error
	InsertLikeInfo(ctx context.Context, biz string, bizId int64, userId int64) error
	DecreaseLikeCount(ctx context.Context, biz string, bizId int64, userId int64) error
	InsertCollectCount(ctx context.Context, biz string, bizId int64, userId int64) error
	DecreaseCollectCount(ctx context.Context, biz string, bizId int64, userId int64) error
	BatchInteractiveReadCount(ctx context.Context, biz string, bizIds, userIds []int64) error
	GetInteractiveInfoByBizId(ctx context.Context, biz string, bizId int64) (domain.InteractiveDomain, error)
}

type InteractiveDaoV1 struct {
	logger loggerx.Logger
	db     *gorm.DB
}

func NewInteractiveDao(db *gorm.DB, l loggerx.Logger) InteractiveDao {
	return &InteractiveDaoV1{
		db:     db,
		logger: l,
	}
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

type InteractiveRead struct {
	Id         int64  `gorm:"primaryKey;autoIncrement"`   // 主键
	Biz        string `gorm:"uniqueIndex:idx_biz_id_biz"` // 业务类型
	BizId      int64  `gorm:"uniqueIndex:idx_biz_id_biz"` // 业务id（文章id）
	UpdateTime int64  // 更新时间
	CreateTime int64  // 创建时间
	UserId     int64  // 用户id
}

func (InteractiveRead) TableName() string {
	return "tb_interactive_read"
}

type InteractiveLike struct {
	Id         int64  `gorm:"primaryKey;autoIncrement"`   // 主键
	Biz        string `gorm:"uniqueIndex:idx_biz_id_biz"` // 业务类型
	BizId      int64  `gorm:"uniqueIndex:idx_biz_id_biz"` // 业务id（文章id）
	UpdateTime int64  // 更新时间
	CreateTime int64  // 创建时间
	Status     uint8  // 软删除，1点赞，2取消点赞
	UserId     int64  // 用户id
}

func (InteractiveLike) TableName() string {
	return "tb_interactive_like"
}

type InteractiveCollect struct {
	Id         int64  `gorm:"primaryKey;autoIncrement"`   // 主键
	Biz        string `gorm:"uniqueIndex:idx_biz_id_biz"` // 业务类型
	BizId      int64  `gorm:"uniqueIndex:idx_biz_id_biz"` // 业务id（文章id）
	UpdateTime int64  // 更新时间
	CreateTime int64  // 创建时间
	Status     uint8  // 软删除，1收藏，2取消收藏
	UserId     int64  // 用户id
}

func (InteractiveCollect) TableName() string {
	return "tb_interactive_collect"
}

type CollectClipDetail struct {
	Id         int64 `gorm:"primaryKey;autoIncrement"` // 主键
	CollectId  int64
	LifeLogId  int64
	CreateTime int64 // 创建时间
	UpdateTime int64 // 更新时间
	Status     uint8
}

func (CollectClipDetail) TableName() string {
	return "tb_collect_clip_detail"
}

// InsertReadCount 增加阅读数，biz业务类型， bizId业务id， userId用户id
func (i *InteractiveDaoV1) InsertReadCount(ctx context.Context, biz string, bizId int64, userId int64) error {
	var interactive Interactive
	var interactiveRead InteractiveRead
	// Gorm闭包事务，gorm帮我们自动控制事务的生命周期（开始，提交，回滚）
	err := i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 初始化interactive
		now := time.Now().UnixMilli()
		interactive.Biz = biz
		interactive.BizId = bizId
		interactive.CreateTime = now
		interactive.UpdateTime = now
		interactive.ReadCount = 1
		// 1.判断用户的阅读记录是否已存在
		readInfoExists := tx.Where("biz_id = ? AND biz = ? AND user_id = ?", bizId, biz, userId).
			First(&interactiveRead).RowsAffected > 0 // 影响记录条数 > 0，表示存在，返回true，影响记录条数 <= 0，表示不存在，返回false
		// 如果是第一次阅读，插入阅读记录
		if !readInfoExists {
			if err := i.InsertReadInfo(ctx, biz, bizId, userId); err != nil {
				i.logger.Error("插入阅读记录失败", loggerx.Error(err),
					loggerx.String("method:", "InteractiveDaoV1:InsertReadCount"))
				return err
			}
		}
		// 2.更新阅读数
		// 	 发生唯一键冲突，就更新数据
		// 	 没有发生唯一键冲突，就插入数据
		err := tx.Where("biz_id = ? AND biz = ?", bizId, biz).
			Clauses(clause.OnConflict{
				DoUpdates: clause.Assignments(map[string]interface{}{
					"read_count":  gorm.Expr("read_count + 1"),
					"update_time": now,
				}),
			}).Create(&interactive).Error
		// 插入阅读数失败
		if err != nil {
			i.logger.Error("增加阅读数失败", loggerx.Error(err),
				loggerx.String("method:", "InteractiveDaoV1:InsertReadCount"))
			return err
		}
		return nil
	})
	if err != nil {
		i.logger.Error("事务执行失败", loggerx.Error(err),
			loggerx.String("method:", "InteractiveDaoV1:InsertReadCount"))
		return err
	}
	return nil
}

// InsertReadInfo 增加阅读记录，biz业务类型， bizId业务id， userId用户id
func (i *InteractiveDaoV1) InsertReadInfo(ctx context.Context, biz string, bizId int64, userId int64) error {
	var interactiveRead InteractiveRead
	now := time.Now().UnixMilli()
	interactiveRead.Biz = biz
	interactiveRead.BizId = bizId
	interactiveRead.CreateTime = now
	interactiveRead.UpdateTime = now
	interactiveRead.UserId = userId
	err := i.db.WithContext(ctx).Create(&interactiveRead).Error
	if err != nil {
		i.logger.Error("增加阅读记录失败", loggerx.Error(err),
			loggerx.String("method:", "InteractiveDaoV1:InsertReadInfo"))
		return err
	}
	return nil
}

// InsertLikeCount 增加点赞数，biz业务类型， bizId业务id， userId用户id
func (i *InteractiveDaoV1) InsertLikeCount(ctx context.Context, biz string, bizId int64, userId int64) error {
	var interactive Interactive
	var interactiveLike InteractiveLike
	// Gorm闭包事务，gorm帮我们自动控制事务的生命周期（开始，提交，回滚）
	err := i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 初始化interactive
		now := time.Now().UnixMilli()
		interactive.Biz = biz
		interactive.BizId = bizId
		interactive.CreateTime = now
		interactive.UpdateTime = now
		interactive.LikeCount = 1
		// 1.判断用户的点赞记录是否已存在
		readInfoExists := tx.Where("biz_id = ? AND biz = ? AND user_id = ?", bizId, biz, userId).
			First(&interactiveLike).
			// 影响记录条数 > 0，表示存在，返回true
			// 影响记录条数 <= 0，表示不存在，返回false
			RowsAffected > 0
		// 如果是第一次点赞，插入点赞记录
		if !readInfoExists {
			if err := i.InsertLikeInfo(ctx, biz, bizId, userId); err != nil {
				i.logger.Error("插入点赞记录失败", loggerx.Error(err),
					loggerx.String("method:", "InteractiveDaoV1:InsertLikeCount"))
				return err
			}
		}
		// 2.更新点赞数
		// 	 发生唯一键冲突，就更新数据
		// 	 没有发生唯一键冲突，就插入数据
		err := tx.Where("biz_id = ? AND biz = ?", bizId, biz).
			Clauses(clause.OnConflict{
				DoUpdates: clause.Assignments(map[string]interface{}{
					"like_count":  gorm.Expr("like_count + 1"),
					"update_time": now,
				}),
			}).Create(&interactive).Error
		// 插入点赞数失败
		if err != nil {
			i.logger.Error("增加点赞数失败", loggerx.Error(err),
				loggerx.String("method:", "InteractiveDaoV1:InsertLikeCount"))
			return err
		}
		return nil
	})
	if err != nil {
		i.logger.Error("事务执行失败", loggerx.Error(err),
			loggerx.String("method:", "InteractiveDaoV1:InsertLikeCount"))
		return err
	}
	return nil
}

// InsertLikeInfo 增加点赞记录，biz业务类型， bizId业务id， userId用户id
func (i *InteractiveDaoV1) InsertLikeInfo(ctx context.Context, biz string, bizId int64, userId int64) error {
	var interactiveLike InteractiveLike
	now := time.Now().UnixMilli()
	interactiveLike.Biz = biz
	interactiveLike.BizId = bizId
	interactiveLike.CreateTime = now
	interactiveLike.UpdateTime = now
	interactiveLike.UserId = userId
	interactiveLike.Status = 1
	err := i.db.WithContext(ctx).Create(&interactiveLike).Error
	if err != nil {
		i.logger.Error("增加点赞记录失败", loggerx.Error(err),
			loggerx.String("method:", "InteractiveDaoV1:InsertReadInfo"))
		return err
	}
	return nil
}

// DecreaseLikeCount 减少点赞数，biz业务类型， bizId业务id， userId用户id
func (i *InteractiveDaoV1) DecreaseLikeCount(ctx context.Context, biz string, bizId int64, userId int64) error {
	now := time.Now().UnixMilli()
	er := i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新点赞数
		err := tx.Model(&Interactive{}).Where("biz = ? AND biz_id = ?", biz, bizId).
			Updates(map[string]interface{}{
				"like_count":  gorm.Expr("like_count - 1"), // 防止 like_count 变为负数
				"update_time": now,
			}).Error
		if err != nil {
			return err
		}
		// 更新点赞记录
		err = tx.Model(&InteractiveLike{}).Where("biz = ? AND biz_id = ? AND user_id = ?", biz, bizId, userId).
			Updates(map[string]interface{}{
				"status":      2,
				"update_time": now,
			}).Error
		if err != nil {
			return err
		}
		return nil
	})
	if er != nil {
		i.logger.Error("事务执行失败", loggerx.Error(er),
			loggerx.String("method:", "InteractiveDaoV1:DecreaseLikeCount"))
		return er
	}
	return nil
}

// InsertCollectInfo 增加收藏记录，biz业务类型， bizId业务id， userId用户id
func (i *InteractiveDaoV1) InsertCollectInfo(ctx context.Context, biz string, bizId, userId int64) error {
	var interactiveCollect InteractiveCollect
	now := time.Now().UnixMilli()
	interactiveCollect.Biz = biz
	interactiveCollect.BizId = bizId
	interactiveCollect.UpdateTime = now
	interactiveCollect.CreateTime = now
	interactiveCollect.Status = 1
	interactiveCollect.UserId = userId
	err := i.db.WithContext(ctx).Create(&interactiveCollect).Error
	if err != nil {
		i.logger.Error("插入收藏记录失败", loggerx.Error(err),
			loggerx.String("method:", "InteractiveDaoV1:InsertCollectInfo"))
		return err
	}
	return nil
}

// DecreaseCollectCount 减少收藏数，biz业务类型， bizId业务id， userId用户id
func (i *InteractiveDaoV1) DecreaseCollectCount(ctx context.Context, biz string, bizId int64, userId int64) error {
	now := time.Now().UnixMilli()
	er := i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新互动表的收藏数，收藏数-1
		err := tx.Where("biz_id = ? AND biz = ?", bizId, biz).Model(&Interactive{}).
			Updates(map[string]any{
				"collect_count": gorm.Expr("collect_count - 1"),
				"update_time":   now,
			}).Error
		if err != nil {
			i.logger.Error("互动表收藏数更新失败", loggerx.Error(err),
				loggerx.String("method:", "InteractiveRepository:DecreaseCollectCount"))
			return err
		}
		// 更新收藏表的状态，状态=2，表示取消收藏
		err = tx.Where("biz_id = ? AND biz = ? AND user_id = ?", bizId, biz, userId).
			Model(&InteractiveCollect{}).Updates(map[string]any{
			"status":      2,
			"update_time": now,
		}).Error
		if err != nil {
			i.logger.Error("收藏表状态更新失败", loggerx.Error(err),
				loggerx.String("method:", "InteractiveRepository:DecreaseCollectCount"))
			return err
		}
		// 更新收藏夹中文章的状态，状态=2，表示取消收藏
		err = tx.Where("life_log_id = ?", bizId).Model(&CollectClipDetail{}).
			Updates(map[string]any{
				"status":      2,
				"update_time": now,
			}).Error
		if err != nil {
			i.logger.Error("文章移除收藏夹失败", loggerx.Error(err),
				loggerx.String("method:", "InteractiveRepository:DecreaseCollectCount"))
			return err
		}
		return nil
	})
	if er != nil {
		i.logger.Error("事务执行失败", loggerx.Error(er),
			loggerx.String("method:", "InteractiveRepository:DecreaseCollectCount"))
		return er
	}
	return nil
}

// InsertCollectCount 增加收藏数，biz业务类型， bizId业务id， userId用户id
func (i *InteractiveDaoV1) InsertCollectCount(ctx context.Context, biz string, bizId int64, userId int64) error {
	var interactive Interactive
	var interactiveCollect InteractiveCollect
	// Gorm闭包事务，gorm帮我们自动控制事务的生命周期（开始，提交，回滚）
	err := i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 初始化interactive
		now := time.Now().UnixMilli()
		interactive.Biz = biz
		interactive.BizId = bizId
		interactive.CreateTime = now
		interactive.UpdateTime = now
		interactive.CollectCount = 1
		// 1.判断用户的收藏记录是否已存在
		readInfoExists := tx.Where("biz_id = ? AND biz = ? AND user_id = ?", bizId, biz, userId).
			First(&interactiveCollect).
			// 影响记录条数 > 0，表示存在，返回true
			// 影响记录条数 <= 0，表示不存在，返回false
			RowsAffected > 0
		// 如果是第一次收藏，插入收藏记录
		if !readInfoExists {
			if err := i.InsertCollectInfo(ctx, biz, bizId, userId); err != nil {
				i.logger.Error("插入收藏记录失败", loggerx.Error(err),
					loggerx.String("method:", "InteractiveDaoV1:InsertCollectCount"))
				return err
			}
		}
		// 2.更新收藏数
		// 	 发生唯一键冲突，就更新数据
		// 	 没有发生唯一键冲突，就插入数据
		err := tx.Where("biz_id = ? AND biz = ?", bizId, biz).
			Clauses(clause.OnConflict{
				DoUpdates: clause.Assignments(map[string]interface{}{
					"collect_count": gorm.Expr("collect_count + 1"),
					"update_time":   now,
				}),
			}).Create(&interactive).Error
		// 插入收藏数失败
		if err != nil {
			i.logger.Error("增加收藏数失败", loggerx.Error(err),
				loggerx.String("method:", "InteractiveDaoV1:InsertCollectCount"))
			return err
		}
		return nil
	})
	if err != nil {
		i.logger.Error("事务执行失败", loggerx.Error(err),
			loggerx.String("method:", "InteractiveDaoV1:InsertCollectCount"))
		return err
	}
	return nil
}

// BatchInteractiveReadCount 批量增加阅读数
func (i *InteractiveDaoV1) BatchInteractiveReadCount(ctx context.Context, biz string,
	bizIds, userIds []int64) error {
	return i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 调用dao层中的其他方法
		txDAO := NewInteractiveDao(tx, i.logger)
		for idx := range bizIds {
			err := txDAO.InsertReadCount(ctx, biz, bizIds[idx], userIds[idx])
			if err != nil {
				// 记录日志即可，因为阅读数少几个，对用户没多少影响
				return err
			}
		}
		return nil
	})
}

// GetInteractiveInfoByBizId 根据文章id获取文章的互动信息
func (i *InteractiveDaoV1) GetInteractiveInfoByBizId(ctx context.Context, biz string, bizId int64) (domain.InteractiveDomain, error) {
	var interactive Interactive
	err := i.db.WithContext(ctx).Where("biz = ? and biz_id = ?", biz, bizId).First(&interactive).Error
	// 有错误
	if err != nil {
		i.logger.Error("获取文章互动信息失败", loggerx.Error(err),
			loggerx.String("method:", "InteractiveDao:GetInteractiveInfoByBizId"))
		return domain.InteractiveDomain{}, err
	}
	// 没有错误
	return domain.InteractiveDomain{
		Id:           interactive.Id,
		Biz:          interactive.Biz,
		BizId:        interactive.BizId,
		ReadCount:    interactive.ReadCount,
		CollectCount: interactive.CollectCount,
		LikeCount:    interactive.LikeCount,
		CreateTime:   interactive.CreateTime,
		UpdateTime:   interactive.UpdateTime,
	}, nil
}
