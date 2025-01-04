package dao

import (
	"context"
	"errors"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"lifelog-grpc/lifeLog/domain"
	"lifelog-grpc/pkg/loggerx"
	"time"
)

type LifeLogDao interface {
	UpdateById(ctx context.Context, lifeLog LifeLog) (domain.LifeLogDomain, error)
	Insert(ctx context.Context, lifeLog LifeLog) (domain.LifeLogDomain, error)
	DeleteByIds(ctx context.Context, ids []int64, public bool) error
	SelectByTitle(ctx context.Context, title string, limit, offset int64) ([]domain.LifeLogDomain, error)
	SelectById(ctx context.Context, id int64, public bool) (domain.LifeLogDomain, error)
	SelectByAuthorId(ctx context.Context, authorId, limit, offset int64) ([]domain.LifeLogDomain, error)
	RevokeById(ctx context.Context, id, authorId int64) error
	// 同步先线上库和制作库的数据
	Sync(ctx context.Context, lifeLog LifeLog) (domain.LifeLogDomain, error)
	Upsert(ctx context.Context, publishLifeLog PublicLifeLog) error
	SelectByIds(ctx context.Context, ids []int64) ([]domain.LifeLogDomain, error)
}

type GormLifeLogDao struct {
	logger loggerx.Logger
	db     *gorm.DB
}

func NewLifeLogDao(db *gorm.DB, l loggerx.Logger) LifeLogDao {
	return &GormLifeLogDao{
		db:     db,
		logger: l,
	}
}

// LifeLog LifeLog表模型
type LifeLog struct {
	Id      int64 `gorm:"primaryKey;autoIncrement"`
	Title   string
	Content string
	// 在作者id和创建时间上建立联合索引
	// 	  查询时：会按照作者id查询，再按照创建时间排序
	//	  在作者id和创建时间添加联合索引，可以减少查询时间
	AuthorId   int64 `gorm:"index=author_id_create_time"`
	CreateTime int64 `gorm:"index=author_id_create_time"`
	UpdateTime int64
	Status     uint8
}

// TableName 设置表名
func (LifeLog) TableName() string {
	return "tb_lifelog"
}

// PublicLifeLog LifeLog表模型
type PublicLifeLog struct {
	Id      int64 `gorm:"primaryKey;autoIncrement"`
	Title   string
	Content string
	// 在作者id和创建时间上建立联合索引
	// 	  查询时：会按照作者id查询，再按照创建时间排序
	//	  在作者id和创建时间添加联合索引，可以减少查询时间
	AuthorId   int64 `gorm:"index"`
	CreateTime int64
	UpdateTime int64
	Status     uint8
}

// TableName 设置表名
func (PublicLifeLog) TableName() string {
	return "tb_publish_lifelog"
}

// 用户模型
type User struct {
	Id            int64 `gorm:"primaryKey;autoIncrement"`
	Email         string
	Password      string
	CreateTime    int64
	UpdateTime    int64
	Phone         string
	WechatUnionId string
	WechatOpenId  string
	NickName      string
}

func (User) TableName() string {
	return "tb_user"
}

// UpdateById 更新LifeLog
func (g *GormLifeLogDao) UpdateById(ctx context.Context, lifeLog LifeLog) (domain.LifeLogDomain, error) {
	lifeLog.UpdateTime = time.Now().UnixMilli()
	// 更新LifeLog
	// where id = ? and author_id = ? ==> 防止跨用户修改
	res := g.db.WithContext(ctx).Where("id = ? and author_id = ?",
		lifeLog.Id, lifeLog.AuthorId).Model(&LifeLog{}).Updates(map[string]any{
		"title":       lifeLog.Title,
		"content":     lifeLog.Content,
		"update_time": lifeLog.UpdateTime,
		"status":      lifeLog.Status,
	})
	// LifeLog更新失败
	if res.Error != nil {
		g.logger.Error("更新LifeLog失败", loggerx.Error(res.Error),
			loggerx.String("method:", "LifeLogDao:UpdateById"))
		return domain.LifeLogDomain{}, res.Error
	}
	// 影响行数为0，说明LifeLog不存在或创建者非法
	if res.RowsAffected == 0 {
		g.logger.Warn("更新LifeLog失败，创作者非法，可能有人在搞你的系统", loggerx.Error(res.Error),
			loggerx.String("method:", "LifeLogDao:UpdateById"))
		return domain.LifeLogDomain{}, errors.New("更新LifeLog失败，创作者非法，可能有人在搞你的系统")
	}
	// LifeLog更新成功
	return domain.LifeLogDomain{
		Id:      lifeLog.Id,
		Title:   lifeLog.Title,
		Content: lifeLog.Content,
		Author: domain.Author{
			Id: lifeLog.AuthorId,
		},
		CreateTime: lifeLog.CreateTime,
		UpdateTime: lifeLog.UpdateTime,
		Status:     lifeLog.Status,
	}, nil
}

// Insert 插入LifeLog
func (g *GormLifeLogDao) Insert(ctx context.Context, lifeLog LifeLog) (domain.LifeLogDomain, error) {
	now := time.Now().UnixMilli()
	lifeLog.CreateTime = now
	lifeLog.UpdateTime = now
	// 插入LifeLog
	err := g.db.WithContext(ctx).Create(&lifeLog).Error
	// LifeLog插入失败
	if err != nil {
		g.logger.Error("插入LifeLog失败", loggerx.Error(err),
			loggerx.String("method:", "LifeLogDao:Insert"))
		return domain.LifeLogDomain{}, err
	}
	// LifeLog插入成功
	return domain.LifeLogDomain{
		Id:      lifeLog.Id,
		Title:   lifeLog.Title,
		Content: lifeLog.Content,
		Author: domain.Author{
			Id: lifeLog.AuthorId,
		},
		CreateTime: lifeLog.CreateTime,
		UpdateTime: lifeLog.UpdateTime,
		Status:     lifeLog.Status,
	}, nil
}

// DeleteByIds 根据id批量删除LifeLog（也可只删除一篇LifeLog）
func (g *GormLifeLogDao) DeleteByIds(ctx context.Context, ids []int64, public bool) error {
	if public == true {
		// 线上库
		err := g.db.WithContext(ctx).Where("id in ?", ids).Delete(&PublicLifeLog{}).Error
		if err != nil {
			g.logger.Error("删除线上库LifeLog失败", loggerx.Error(err),
				loggerx.String("method:", "LifeLogDao:DeleteById"))
			return err
		}
		return nil
	}
	// 制作库
	err := g.db.WithContext(ctx).Where("id in ?", ids).Delete(&LifeLog{}).Error
	if err != nil {
		g.logger.Error("删除制作库LifeLog失败", loggerx.Error(err),
			loggerx.String("method:", "LifeLogDao:DeleteById"))
		return err
	}
	return nil
}

// SelectByTitle 根据标题查询LifeLog（线上库）
func (g *GormLifeLogDao) SelectByTitle(ctx context.Context, title string, limit, offset int64) ([]domain.LifeLogDomain, error) {
	var ad []domain.LifeLogDomain
	var publicLifeLogs []PublicLifeLog
	return ad, g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 根据LifeLogtitle查询线上库中的LifeLog，模糊查询
		err := g.db.WithContext(ctx).Where("title like ?", "%"+title+"%").
			Limit(int(limit)).
			Offset(int(offset)).
			Order("update_time desc").
			Find(&publicLifeLogs).Error
		// 查询失败
		if err != nil {
			g.logger.Error("查询LifeLog失败", loggerx.Error(err),
				loggerx.String("method:", "LifeLogDao:SelectByTitle"))
			ad = nil
			return err
		}
		// 查找LifeLog对应的作者信息
		// 遍历LifeLog列表，获取作者id
		var authorIds []int64
		var user User
		var userNames []string
		for _, lifeLog := range publicLifeLogs {
			authorIds = append(authorIds, lifeLog.AuthorId)
		}
		for _, id := range authorIds {
			err = tx.WithContext(ctx).Where("id = ?", id).Find(&user).Error
			userNames = append(userNames, user.NickName)
		}
		if err != nil {
			g.logger.Error("查询作者信息失败", loggerx.Error(err),
				loggerx.String("method:", "LifeLogDao:SelectByTitle"))
			ad = nil
			return err
		}
		// 查询成功
		atads := g.lifeLogsToPublicLifeLogDomains(publicLifeLogs)
		ad = g.usersToLifeLogDomains(userNames, atads)
		return nil
	})
}

// SelectById 根据id查询LifeLog
func (g *GormLifeLogDao) SelectById(ctx context.Context, id int64, public bool) (domain.LifeLogDomain, error) {
	var lifeLog LifeLog
	var publicLifeLog PublicLifeLog
	// 线上库
	if public == true {
		err := g.db.WithContext(ctx).Where("id = ? and status = ?",
			id, domain.LifeLogStatusPublished).Find(&publicLifeLog).Error
		if err != nil {
			g.logger.Error("查询线上库LifeLog失败", loggerx.Error(err),
				loggerx.String("method:", "LifeLogDao:SelectById"))
			return domain.LifeLogDomain{}, err
		}
		return domain.LifeLogDomain{
			Id:      publicLifeLog.Id,
			Title:   publicLifeLog.Title,
			Content: publicLifeLog.Content,
			Author: domain.Author{
				Id: publicLifeLog.AuthorId,
			},
			CreateTime: publicLifeLog.CreateTime,
			UpdateTime: publicLifeLog.UpdateTime,
			Status:     publicLifeLog.Status,
		}, nil
	}
	// 制作库
	err := g.db.WithContext(ctx).Where("id = ?", id).Find(&lifeLog).Error
	if err != nil {
		g.logger.Error("查询制作库LifeLog失败", loggerx.Error(err),
			loggerx.String("method:", "LifeLogDao:SelectById"))
		return domain.LifeLogDomain{}, err
	}
	return domain.LifeLogDomain{
		Id:      lifeLog.Id,
		Title:   lifeLog.Title,
		Content: lifeLog.Content,
		Author: domain.Author{
			Id: lifeLog.AuthorId,
		},
		CreateTime: lifeLog.CreateTime,
		UpdateTime: lifeLog.UpdateTime,
		Status:     lifeLog.Status,
	}, nil
}

// SelectByIds 根据id批量查询LifeLog
func (g *GormLifeLogDao) SelectByIds(ctx context.Context, ids []int64) ([]domain.LifeLogDomain, error) {
	var publicLifeLog []PublicLifeLog
	// 线上库
	err := g.db.WithContext(ctx).Where("ids in ? and status = ?",
		ids, domain.LifeLogStatusPublished).Find(&publicLifeLog).Error
	if err != nil {
		g.logger.Error("查询线上库LifeLog失败", loggerx.Error(err),
			loggerx.String("method:", "LifeLogDao:SelectByIds"))
		return nil, err
	}
	// 将publicLifeLog中的数据复制到lifelogs中
	var lifeLogs []domain.LifeLogDomain
	copier.Copy(&lifeLogs, &publicLifeLog)
	return lifeLogs, nil
}

// SelectByAuthorId 根据作者id查询LifeLog
func (g *GormLifeLogDao) SelectByAuthorId(ctx context.Context, authorId, limit, offset int64) ([]domain.LifeLogDomain, error) {
	var lifeLogs []LifeLog
	var ad []domain.LifeLogDomain
	return ad, g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 分页查询
		err := g.db.WithContext(ctx).Where("author_Id = ?", authorId).
			Limit(int(limit)).
			Offset(int(offset)).
			Order("update_time desc").
			Find(&lifeLogs).Error
		// 查询失败
		if err != nil {
			g.logger.Error("查询创作者，制作库LifeLog列表失败", loggerx.Error(err),
				loggerx.String("method:", "LifeLogDao:SelectByAuthorId"))
			ad = nil
			return err
		}
		// 查找LifeLog对应的作者信息
		// 遍历LifeLog列表，获取作者id
		var authorIds []int64
		var user User
		var userNames []string
		for _, lifeLog := range lifeLogs {
			authorIds = append(authorIds, lifeLog.AuthorId)
		}
		for _, id := range authorIds {
			err = tx.WithContext(ctx).Where("id = ?", id).Find(&user).Error
			userNames = append(userNames, user.NickName)
		}
		if err != nil {
			g.logger.Error("查询作者信息失败", loggerx.Error(err),
				loggerx.String("method:", "LifeLogDao:SelectByAuthorId"))
			ad = nil
			return err
		}
		// 查询成功
		atads := g.lifeLogsToLifeLogDomains(lifeLogs)
		ad = g.usersToLifeLogDomains(userNames, atads)
		return nil
	})
}

// RevokeById 撤销LifeLog
func (g *GormLifeLogDao) RevokeById(ctx context.Context, id, authorId int64) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新线先上库的status字段
		res := g.db.WithContext(ctx).Model(&PublicLifeLog{}).
			// 保证只有作者本人可以撤销
			Where("id = ? and author_id = ?", id, authorId).
			Updates(map[string]any{
				"status":      domain.LifeLogStatusHidden,
				"update_time": now,
			})
		if res.Error != nil {
			g.logger.Error("撤销LifeLog失败，更新线上库的status字段错误", loggerx.Error(res.Error),
				loggerx.String("method:", "LifeLogDao:RevokeById"))
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("没有找到LifeLog，可能有人在搞你的系统")
		}
		// 更新制作库的status字段
		res = g.db.WithContext(ctx).Model(&LifeLog{}).
			Where("id = ? and author_id = ?", id, authorId).
			Updates(map[string]any{
				"status":      domain.LifeLogStatusHidden,
				"update_time": now,
			})
		if res.Error != nil {
			g.logger.Error("撤销LifeLog失败，更新制作库的status字段错误", loggerx.Error(res.Error),
				loggerx.String("method:", "LifeLogDao:RevokeById"))
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("没有找到LifeLog，可能有人在搞你的系统")
		}
		return nil
	})
}

// Sync 同步线上库和制作库的数据
func (g *GormLifeLogDao) Sync(ctx context.Context, lifeLog LifeLog) (domain.LifeLogDomain, error) {
	// 先操作制作库(数据库表)，再操作线上库(数据库表)
	var id = lifeLog.Id
	// 开启事务，采用闭包形态，gorm自动完成，begin，Rollback，commit操作，不需要我们操心
	//    tx <==> Transaction事务
	err := g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		// 操作制作库（复用已有的数据库事务，而不是创建一个新的数据库连接）
		txDao := NewLifeLogDao(tx, g.logger)
		// 传入了id，表示是更新数据库
		if id > 0 {
			_, err = txDao.UpdateById(ctx, lifeLog)
		} else {
			// 没传入id，表示插入数据
			// 这是属于，直接发布lifelog，制作库的lifelog的status=1
			lifeLog.Status = domain.LifeLogStatusPublished
			_, err = txDao.Insert(ctx, lifeLog)
		}
		if err != nil {
			g.logger.Error("操作制作库失败", loggerx.Error(err),
				loggerx.String("method:", "LifeLogDao:Sync"))
			return err
		}
		// 组装线上库数据
		publishLifeLog := PublicLifeLog{
			Id:         id,
			Title:      lifeLog.Title,
			Content:    lifeLog.Content,
			AuthorId:   lifeLog.AuthorId,
			CreateTime: lifeLog.CreateTime,
			UpdateTime: lifeLog.UpdateTime,
			Status:     uint8(domain.LifeLogStatusPublished),
		}
		// 插入或更新线上库
		return txDao.Upsert(ctx, publishLifeLog)
	})
	return domain.LifeLogDomain{
		Id:      id,
		Title:   lifeLog.Title,
		Content: lifeLog.Content,
		Author: domain.Author{
			Id: lifeLog.AuthorId,
		},
		CreateTime: lifeLog.CreateTime,
		UpdateTime: lifeLog.UpdateTime,
	}, err
}

// Upsert 插入或更新线上库
func (a *GormLifeLogDao) Upsert(ctx context.Context, publishLifeLog PublicLifeLog) error {
	now := time.Now().UnixMilli()
	publishLifeLog.UpdateTime = now
	publishLifeLog.CreateTime = now
	// 最终生成的mysql语句为：insert xxx on duplicate key update xxx
	//    插入数据时，如果发生了冲突(数据库中已经存在数据)，就更新数据
	err := a.db.WithContext(ctx).Clauses(clause.OnConflict{ // OnConflict数据冲突
		// 要更新数据
		DoUpdates: clause.Assignments(map[string]any{
			"title":       publishLifeLog.Title,
			"content":     publishLifeLog.Content,
			"update_time": now,
			"status":      publishLifeLog.Status,
		}),
		// 要插入数据
	}).Create(&publishLifeLog).Error
	if err != nil {
		return err // 返回错误
	}
	return nil
}

// lifeLogsToLifeLogDomains 将制作库LifeLog列表转换为领域模型
func (g *GormLifeLogDao) lifeLogsToLifeLogDomains(lifeLogs []LifeLog) []domain.LifeLogDomain {
	ads := make([]domain.LifeLogDomain, 0, len(lifeLogs))
	for _, lifeLog := range lifeLogs {
		ads = append(ads, domain.LifeLogDomain{
			Id:      lifeLog.Id,
			Title:   lifeLog.Title,
			Content: lifeLog.Content,
			Author: domain.Author{
				Id: lifeLog.AuthorId,
			},
			CreateTime: lifeLog.CreateTime,
			UpdateTime: lifeLog.UpdateTime,
		})
	}
	return ads
}

// 将Users信息插入到LifeLogDomains中
func (g *GormLifeLogDao) usersToLifeLogDomains(userNames []string, atads []domain.LifeLogDomain) []domain.LifeLogDomain {
	for i, userName := range userNames {
		atads[i].Author.Name = userName
	}
	return atads
}

// lifeLogsToLifeLogDomains 将线上库LifeLog列表转换为领域模型
func (g *GormLifeLogDao) lifeLogsToPublicLifeLogDomains(lifeLogs []PublicLifeLog) []domain.LifeLogDomain {
	ads := make([]domain.LifeLogDomain, 0, len(lifeLogs))
	for _, lifeLog := range lifeLogs {
		ads = append(ads, domain.LifeLogDomain{
			Id:      lifeLog.Id,
			Title:   lifeLog.Title,
			Content: lifeLog.Content,
			Author: domain.Author{
				Id: lifeLog.AuthorId,
			},
			CreateTime: lifeLog.CreateTime,
			UpdateTime: lifeLog.UpdateTime,
		})
	}
	return ads
}
