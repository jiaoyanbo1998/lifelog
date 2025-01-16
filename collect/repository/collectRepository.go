package repository

import (
	"context"
	"lifelog-grpc/collect/domain"
	"lifelog-grpc/collect/repository/cache"
	"lifelog-grpc/collect/repository/dao"
)

type CollectRepository interface {
	SaveCollect(ctx context.Context, collectDomain domain.CollectDomain) error
	DeleteCollect(ctx context.Context, ids []int64, authorId int64) error
	CollectList(ctx context.Context, userId int64, limit int, offset int) ([]domain.CollectDomain, error)
	InsertCollectDetail(ctx context.Context, detailDomain domain.CollectDetailDomain) error
	CollectDetail(ctx context.Context, collectId int64, limit int, offset int, authorId int64) ([]domain.CollectDetailDomain, error)
	DeleteCollectDetail(ctx context.Context, collectId, lifeLogId, authorId int64) error
}

type CollectRepositoryV1 struct {
	collectDao   dao.CollectDao
	collectCache cache.CollectCache
}

func NewCollectRepository(collectDao dao.CollectDao, collectCache cache.CollectCache) CollectRepository {
	return &CollectRepositoryV1{
		collectDao:   collectDao,
		collectCache: collectCache,
	}
}

// DeleteCollectDetail 删除收藏夹详情
func (c *CollectRepositoryV1) DeleteCollectDetail(ctx context.Context, collectId, lifeLogId, authorId int64) error {
	return c.collectDao.DeleteCollectDetail(ctx, collectId, lifeLogId, authorId)
}

// SaveCollect 编辑收藏夹
func (c *CollectRepositoryV1) SaveCollect(ctx context.Context, collectDomain domain.CollectDomain) error {
	// 更新收藏夹
	if collectDomain.Id > 0 {
		return c.collectDao.UpdateCollect(ctx, dao.Collect{
			Id:       collectDomain.Id,
			Name:     collectDomain.Name,
			AuthorId: collectDomain.AuthorId,
		})
	}
	return c.collectDao.InsertCollect(ctx, dao.Collect{
		Name:     collectDomain.Name,
		AuthorId: collectDomain.AuthorId,
	})
}

// DeleteCollect 删除收藏夹
func (c *CollectRepositoryV1) DeleteCollect(ctx context.Context, ids []int64, authorId int64) error {
	return c.collectDao.DeleteCollectByIds(ctx, ids, authorId)
}

// CollectList 收藏夹列表
func (c *CollectRepositoryV1) CollectList(ctx context.Context, userId int64, limit int, offset int) ([]domain.CollectDomain, error) {
	// 从缓存中获取
	res, err := c.collectCache.Get(ctx, userId)
	// 缓存中有，直接返回
	if err == nil {
		return res, nil
	}
	// 缓存中没有，从数据库中获取
	res, err = c.collectDao.PageQuery(ctx, userId, limit, offset)
	// 数据库查询失败
	if err != nil {
		return nil, err
	}
	// 回写缓存
	err = c.collectCache.Set(ctx, userId, res)
	// 缓存回写失败
	if err != nil {
		return nil, err
	}
	// 数据库查询成功，返回结果
	return res, nil
}

// InsertCollectDetail 将文章插入到收藏夹
func (c *CollectRepositoryV1) InsertCollectDetail(ctx context.Context, detailDomain domain.CollectDetailDomain) error {
	return c.collectDao.InsertCollectDetail(ctx, dao.CollectDetail{
		LifeLogId: detailDomain.LifeLogId,
		CollectId: detailDomain.CollectId,
		AuthorId:  detailDomain.AuthorId,
	})
}

// CollectDetail 收藏夹详情
func (c *CollectRepositoryV1) CollectDetail(ctx context.Context, collectId int64, limit int, offset int,
	authorId int64) ([]domain.CollectDetailDomain, error) {
	return c.collectDao.GetCollectDetailById(ctx, collectId, limit, offset, authorId)
}
