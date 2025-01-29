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
	SearchByIds(ctx context.Context, ids []int64) ([]domain.LifeLogDomain, error)
}

type LifeLogRepositoryV1 struct {
	lifeLogDao   dao.LifeLogDao
	lifeLogCache cache.LifeLogCache
	localCache   *cache.LocalCache
	logger       loggerx.Logger
}

func NewLifeLogRepository(lifeLogDao dao.LifeLogDao, l loggerx.Logger,
	lifeLogCache cache.LifeLogCache, localCache *cache.LocalCache) LifeLogRepository {
	return &LifeLogRepositoryV1{
		lifeLogDao:   lifeLogDao,
		lifeLogCache: lifeLogCache,
		logger:       l,
		localCache:   localCache,
	}
}

// Modify 修改LifeLog
func (a *LifeLogRepositoryV1) Modify(ctx context.Context, lifeLogDomain domain.LifeLogDomain) (domain.LifeLogDomain, error) {
	go func() {
		// 删除本地缓存
		err := a.localCache.DelFirstPage(lifeLogDomain.Author.Id)
		if err != nil {
			a.logger.Error("删除本地缓存失败", loggerx.Error(err),
				loggerx.String("method:", "lifeLogRepository:Modify"))
		}
		// 删除redis缓存
		err = a.lifeLogCache.DelFirstPage(ctx, lifeLogDomain.Author.Id)
		if err != nil {
			a.logger.Error("删除redis缓存失败", loggerx.Error(err),
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

// Create 创建LifeLog（制作库）
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
		Status:   domain.LifeLogStatusUnPublish,
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
	// 查询本地缓存
	detail, err := a.localCache.Get(id)
	if err == nil {
		return detail, nil
	}
	// 查询redis缓存
	detail, err = a.lifeLogCache.GetFirstPageDetail(ctx, id)
	// 缓存命中，直接返回
	if err == nil {
		return detail, nil
	}
	// 缓存未命中，查询数据库
	return a.lifeLogDao.SelectById(ctx, id, public)
}

// SearchByIds 搜索多个LifeLog
func (a *LifeLogRepositoryV1) SearchByIds(ctx context.Context, ids []int64) ([]domain.LifeLogDomain, error) {
	return a.lifeLogDao.SelectByIds(ctx, ids)
}

// SearchByAuthorId 搜索LifeLog
func (a *LifeLogRepositoryV1) SearchByAuthorId(ctx context.Context, authorId, limit, offset int64) ([]domain.LifeLogDomain, error) {
	// 参数校验
	if authorId <= 0 || limit < 0 || offset < 0 {
		return nil, errors.New("invalid parameters")
	}
	// 查询本地缓存
	if isFirstPage(offset, limit) {
		ads, err := a.localCache.GetFirstPage(authorId)
		if err == nil && len(ads) > 0 {
			return ads[:min(limit, int64(len(ads)))], nil
		}
	}
	// 本地缓存没有数据，查询redis缓存
	if isFirstPage(offset, limit) {
		ads, err := a.lifeLogCache.GetFirstPage(ctx, authorId)
		if err == nil && len(ads) > 0 {
			return ads[:min(limit, int64(len(ads)))], nil
		}
	}
	// redis没有缓存，查询数据库
	ads, err := a.lifeLogDao.SelectByAuthorId(ctx, authorId, limit, offset)
	if err != nil {
		a.logger.Error("查询数据库失败", loggerx.Error(err),
			loggerx.String("method:", "lifeLogRepository:SearchByAuthorId"),
			loggerx.Int64("authorId", authorId),
			loggerx.Int64("limit", limit),
			loggerx.Int64("offset", offset))
		// 所有数据源都查询失败
		// 返回降级提示
		return nil, errors.New("系统繁忙，请稍后再试")
	}
	// 异步回写缓存
	if isFirstPage(offset, limit) {
		if err := a.localCache.SetFirstPage(authorId, ads); err != nil {
			a.logger.Warn("回写本地缓存失败", loggerx.Error(err),
				loggerx.String("method:", "lifeLogRepository:SearchByAuthorId"),
				loggerx.Int64("authorId", authorId))
		}
		if err := a.lifeLogCache.SetFirstPage(ctx, authorId, ads); err != nil {
			a.logger.Warn("回写redis缓存失败", loggerx.Error(err),
				loggerx.String("method:", "lifeLogRepository:SearchByAuthorId"),
				loggerx.Int64("authorId", authorId))
		}
	}
	if err := a.preCache(ctx, ads); err != nil {
		a.logger.Warn("预加载缓存失败", loggerx.Error(err),
			loggerx.String("method:", "lifeLogRepository:SearchByAuthorId"),
			loggerx.Int64("authorId", authorId))
	}
	return ads, nil
}

func isFirstPage(offset, limit int64) bool {
	return offset == 0 && limit <= 100
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// preCache 预缓存第一页的第一条数据
func (a *LifeLogRepositoryV1) preCache(ctx context.Context, ads []domain.LifeLogDomain) error {
	if len(ads) > 0 {
		// 本地缓存
		err := a.localCache.Set(ads[0])
		if err != nil {
			return err
		}
		// redis
		err = a.lifeLogCache.Set(ctx, ads[0])
		if err != nil {
			return err
		}
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
