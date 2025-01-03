package repository

import (
	"context"
	"errors"
	"lifelog-grpc/lifeLog/domain"
	"lifelog-grpc/lifeLog/repository/cache"
	"lifelog-grpc/lifeLog/repository/dao"
	"lifelog-grpc/pkg/loggerx"
)

type LifeLogRepository interface {
	Modify(ctx context.Context, lifeLogDomain domain.LifeLogDomain) (domain.LifeLogDomain, error)
	Create(ctx context.Context, lifeLogDomain domain.LifeLogDomain) (domain.LifeLogDomain, error)
	Delete(ctx context.Context, ids []int64, public bool) error
	SearchByTitle(ctx context.Context, title string, limit, offset int64) ([]domain.LifeLogDomain, error)
	SearchById(ctx context.Context, id int64, public bool) (domain.LifeLogDomain, error)
	SearchByAuthorId(ctx context.Context, authorId, limit, offset int64) ([]domain.LifeLogDomain, error)
	RevokeById(ctx context.Context, id, authorId int64) error
	Sync(ctx context.Context, lifeLogDomain domain.LifeLogDomain) error
}

type LifeLogRepositoryV1 struct {
	lifeLogDao   dao.LifeLogDao
	lifeLogCache cache.LifeLogCache
	logger       loggerx.Logger
}

func NewLifeLogRepository(lifeLogDao dao.LifeLogDao, l loggerx.Logger,
	lifeLogCache cache.LifeLogCache) LifeLogRepository {
	return &LifeLogRepositoryV1{
		lifeLogDao:   lifeLogDao,
		lifeLogCache: lifeLogCache,
		logger:       l,
	}
}

// Modify 修改LifeLog
func (a *LifeLogRepositoryV1) Modify(ctx context.Context, lifeLogDomain domain.LifeLogDomain) (domain.LifeLogDomain, error) {
	go func() {
		err := a.lifeLogCache.DelFirstPage(ctx, lifeLogDomain.Author.Id)
		if err != nil {
			a.logger.Error("删除缓存失败", loggerx.Error(err),
				loggerx.String("method:", "lifeLogRepository:Modify"))
		}
	}()
	return a.lifeLogDao.UpdateById(ctx, dao.LifeLog{
		Id:       lifeLogDomain.Id,
		Title:    lifeLogDomain.Title,
		Content:  lifeLogDomain.Content,
		AuthorId: lifeLogDomain.Author.Id,
		// 制作库中的LifeLog，只有未发布状态
		Status: domain.LifeLogStatusUnPublish,
	})
}

// Create 创建LifeLog
func (a *LifeLogRepositoryV1) Create(ctx context.Context, lifeLogDomain domain.LifeLogDomain) (domain.LifeLogDomain, error) {
	go func() {
		err := a.lifeLogCache.DelFirstPage(ctx, lifeLogDomain.Author.Id)
		if err != nil {
			a.logger.Error("删除缓存失败", loggerx.Error(err),
				loggerx.String("method:", "lifeLogRepository:Create"))
		}
	}()
	return a.lifeLogDao.Insert(ctx, dao.LifeLog{
		Title:    lifeLogDomain.Title,
		Content:  lifeLogDomain.Content,
		AuthorId: lifeLogDomain.Author.Id,
		// 制作库中的LifeLog，只有未发布状态
		Status: domain.LifeLogStatusUnPublish,
	})
}

// Delete 删除LifeLog
func (a *LifeLogRepositoryV1) Delete(ctx context.Context, ids []int64, public bool) error {
	return a.lifeLogDao.DeleteByIds(ctx, ids, public)
}

// SearchByTitle 搜索LifeLog
func (a *LifeLogRepositoryV1) SearchByTitle(ctx context.Context, title string, limit, offset int64) ([]domain.LifeLogDomain, error) {
	return a.lifeLogDao.SelectByTitle(ctx, title, limit, offset)
}

// SearchById 搜索LifeLog
func (a *LifeLogRepositoryV1) SearchById(ctx context.Context, id int64, public bool) (domain.LifeLogDomain, error) {
	return a.lifeLogDao.SelectById(ctx, id, public)
}

// SearchByAuthorId 搜索LifeLog
func (a *LifeLogRepositoryV1) SearchByAuthorId(ctx context.Context, authorId, limit, offset int64) ([]domain.LifeLogDomain, error) {
	// 缓存第一页数据（我不知道第一页有多少条数据，暂且认为有100条）
	if offset == 0 && limit <= 100 {
		ads, err := a.lifeLogCache.GetFirstPage(ctx, authorId)
		if err == nil && len(ads) > 0 {
			return ads[:limit], nil
		}
	}
	// 根据作者id，查询作者的LifeLog列表
	ads, err := a.lifeLogDao.SelectByAuthorId(ctx, authorId, limit, offset)
	if err != nil {
		a.logger.Error("查询数据库失败", loggerx.Error(err),
			loggerx.String("method:", "lifeLogRepository:SearchByAuthorId"))
		return nil, err
	}
	// 将查询出来的数据，存储到redis中
	err = a.lifeLogCache.SetFirstPage(ctx, authorId, ads)
	if err != nil {
		a.logger.Error("回写缓存失败", loggerx.Error(err),
			loggerx.String("method:", "lifeLogRepository:SearchByAuthorId"))
		return ads, err
	}
	// 缓存第一页的第一条数据
	err = a.preCache(ctx, ads)
	if err != nil {
		a.logger.Error("预缓存失败", loggerx.Error(err),
			loggerx.String("method:", "lifeLogRepository:SearchByAuthorId"))
	}
	return ads, nil
}

// preCache 预缓存第一页的第一条数据
func (a *LifeLogRepositoryV1) preCache(ctx context.Context, ads []domain.LifeLogDomain) error {
	if len(ads) > 0 {
		return a.lifeLogCache.Set(ctx, ads[0])
	}
	return errors.New("没有数据")
}

// RevokeById 撤销LifeLog
func (a *LifeLogRepositoryV1) RevokeById(ctx context.Context, id, authorId int64) error {
	return a.lifeLogDao.RevokeById(ctx, id, authorId)
}

// Sync 同步LifeLog
func (a *LifeLogRepositoryV1) Sync(ctx context.Context, lifeLogDomain domain.LifeLogDomain) error {
	// 防止缓存不一致，在数据同步前，将缓存中的数据删除
	go func() {
		err := a.lifeLogCache.DelFirstPage(ctx, lifeLogDomain.Author.Id)
		if err != nil {
			a.logger.Error("删除缓存失败", loggerx.Error(err),
				loggerx.String("method:", "lifeLogRepository:Modify"))
		}
	}()
	// 同步数据
	ad, err := a.lifeLogDao.Sync(ctx, dao.LifeLog{
		Id:       lifeLogDomain.Id,
		Title:    lifeLogDomain.Title,
		Content:  lifeLogDomain.Content,
		AuthorId: lifeLogDomain.Author.Id,
		Status:   lifeLogDomain.Status,
	})
	// 数据同步没有出错
	if err == nil {
		// 将第一次发布的LifeLog缓存到redis中
		er := a.lifeLogCache.SetPublic(ctx, ad)
		if er != nil {
			a.logger.Warn("缓存失败", loggerx.Error(err),
				loggerx.String("method:", "lifeLogRepository:Sync"))
		}
	}
	// 数据同步出错
	return err
}
